package teleinfo

import (
	"bufio"
	"bytes"
	"fmt"
)

// Read return data from Teleinfo
// Example: PAPP:00710 MOTDETAT:000000 ADCO:0123456789012 OPTARIF:BASE BASE:012345678 PTEC:TH.. ISOUSC:30 IINST:003 IMAX:027
func (c *Client) Read() (map[string]string, error) {
	rawFrame, err := readRawFrame(c.buffer)
	if err != nil {
		return nil, err
	}
	return decodeFrame(rawFrame)
}

func decodeFrame(rawFrame []byte) (map[string]string, error) {
	const (
		checksumLength = 1
	)

	strFrame := bytes.Trim(rawFrame, "\r\n")

	fields := bytes.Split(strFrame, []byte("\r\n"))
	info := make(map[string]string)
	for _, field := range fields {
		elts := bytes.SplitN(field, []byte(" "), 3)

		if len(elts) != 3 {
			return nil, fmt.Errorf("error decoding frame, invalid number of elements for data (data: '%s')", field)
		}
		name, value, trail := elts[0], elts[1], elts[2]

		if len(trail) != checksumLength {
			return nil, fmt.Errorf("error decoding frame, invalid checksum length (actual: %d, expected: %d)", len(trail), checksumLength)
		}
		readChecksum := byte(trail[0])
		expectedChecksum := computeChecksum(name, value)
		if readChecksum != expectedChecksum {
			return nil, fmt.Errorf("error decoding frame, invalid checksum (field: '%s', value: '%s', read: '%c', expected: '%c'", name, value, readChecksum, expectedChecksum)
		}
		info[string(name)] = string(value)
	}
	return info, nil
}

func sum(a []byte) (res byte) {
	res = 0
	for _, c := range a {
		res += c
	}
	return
}

func computeChecksum(name []byte, value []byte) byte {
	// NOTE: 0x20 == ASCII space char
	checksum := sum(name) + byte(0x20) + sum(value)

	// Map to a single char E [0x20;0x7F]
	checksum = (checksum & 0x3F) + 0x20
	return checksum
}

func readRawFrame(buffer *bufio.Reader) ([]byte, error) {
	const (
		FrameStart byte = 0x2
		FrameEnd   byte = 0x3
	)
	var frame []byte

	// TODO: check for overflow
	// TODO: check for interrupted frame marker
	if _, err := buffer.ReadSlice(FrameStart); err != nil {
		return nil, fmt.Errorf("error looking for start of frame marker (%s)", err)
	}

	var errRead error
	if frame, errRead = buffer.ReadBytes(FrameEnd); errRead != nil {
		return nil, fmt.Errorf("error looking for end of frame marker (%s)", errRead)
	}

	if len(frame) == 0 {
		return frame, fmt.Errorf("read empty frame")
	}
	frame = frame[0 : len(frame)-1]
	return frame, nil
}
