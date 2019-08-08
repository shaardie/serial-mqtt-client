package main

import (
	"io"

	"github.com/tarm/serial"
)

func newSerial(port string, baudrate int) (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: port, Baud: baudrate}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return s, nil
}
