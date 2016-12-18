package main

import (
	"fmt"
	"os"

	"github.com/yesnault/teleinfo"
)

func main() {
	c, err := teleinfo.NewClient(teleinfo.Options{Device: "/dev/ttyUSB0"})
	if err != nil {
		fmt.Printf("Error while initialize teleinfo client:%s\n", err)
		os.Exit(1)
	}

	defer c.Close()

	for {
		info, err := c.Read()
		if err != nil {
			fmt.Printf("Error reading (%s)\n", err)
			os.Exit(1)
		}
		fmt.Printf("Info: %+v\n", info)
	}

}
