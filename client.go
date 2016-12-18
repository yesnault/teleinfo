package teleinfo

import (
	"bufio"
	"fmt"

	"github.com/tarm/serial"
)

// Client represents a Client for Teleinfo
type Client struct {
	device string
	port   *serial.Port
	buffer *bufio.Reader
}

//Options is a struct to initialize a Teleinfo client
type Options struct {
	Device string
}

// Open open where teleinfo is connected
// example: /dev/ttyUSB0
func (c *Client) open() error {
	cfg := &serial.Config{
		Name:     c.device,
		Baud:     1200,
		Size:     7,
		Parity:   serial.ParityEven,
		StopBits: serial.Stop1,
	}
	var err error

	c.port, err = serial.OpenPort(cfg)
	if err == nil {
		c.buffer = bufio.NewReader(c.port)
	}
	return err
}

// Close close port of Teleinfo Client
func (c *Client) Close() {
	c.port.Close()
}

//NewClient initialize a Teleinfo client
func NewClient(opts Options) (*Client, error) {
	if opts.Device == "" {
		return nil, fmt.Errorf("Invalid configuration, please Device")
	}
	c := &Client{
		device: opts.Device,
	}
	return c, c.open()
}
