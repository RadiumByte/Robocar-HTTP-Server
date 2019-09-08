package ral

import "fmt"

// RoboCar represents Raspberry Pi based car
type RoboCar struct {
}

// SendCommand creates HTTP client and sends coommand to robot
func (robot *RoboCar) SendCommand(command string) {
	fmt.Println("Sending command: " + command)
}

// NewRoboCar constructs object of RoboCar
func NewRoboCar() (*RoboCar, error) {
	res := &RoboCar{}

	return res, nil
}
