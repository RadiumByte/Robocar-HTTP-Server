package ral

import (
	"fmt"
	"github.com/tarm/serial"
	"strconv"
	"time"
)

// RoboCar represents Raspberry Pi based car
type RoboCar struct {
	SerialConfig  *serial.Config
	SerialPort    *serial.Port
	CurrentPort   int
}

func (robot *RoboCar) Watchdog() {
	fmt.Println("Car is initializing")
	time.Sleep(20 * time.Second)
	robot.SendCommand("START")

	for {
		robot.SendCommand("ACK")
		time.Sleep(200 * time.Millisecond)
	}
}

// ConnectToChassis finds Car's Serial port by scanning and connects to it
func (robot *RoboCar) ConnectToChassis() error {
	var err error
	robot.CurrentPort = 0
	for robot.CurrentPort < 20 {
		robot.SerialConfig = &serial.Config{Name: "/dev/ttyACM" + strconv.Itoa(robot.CurrentPort), Baud: 38400}
		robot.SerialPort, err = serial.OpenPort(robot.SerialConfig)
		if err == nil {
			return nil
		}
		robot.CurrentPort += 1
	}
	robot.CurrentPort = 0
	return err
}

// SendCommand creates HTTP client and sends coommand to robot
func (robot *RoboCar) SendCommand(command string) {
	fmt.Println("Sending command to car: " + command)
	command += "\r\n"
	_, err := robot.SerialPort.Write([]byte(command))
	if err != nil {
		fmt.Println("Error while sending a command to car")
		err = robot.ConnectToChassis()

		if err != nil {
			fmt.Println("No communication with chassis")
			panic("Failure")
		}
		_, err = robot.SerialPort.Write([]byte(command))
	}
}

// NewRoboCar constructs object of RoboCar
func NewRoboCar() (*RoboCar, error) {
	res := &RoboCar{}

	res.CurrentPort = 0
	err := res.ConnectToChassis()

	if err != nil {
		fmt.Println("No communication with chassis")
		panic("error")
	}
	go res.Watchdog()

	return res, nil
}
