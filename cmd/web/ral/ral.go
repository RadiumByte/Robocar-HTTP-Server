package ral

import "fmt"

//"github.com/valyala/fasthttp"

// RoboCar represents Raspberry Pi based car
type RoboCar struct {
	CarIP   string
	CarPort string
}

// SendCommand creates HTTP client and sends coommand to robot
func (robot *RoboCar) SendCommand(command string) {
	fmt.Println("Sending command: " + command)
	//url := "http://" + robot.CarIP + robot.CarPort + "/" + command
	/*
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(url)
		req.Header.SetMethod("PUT")

		resp := fasthttp.AcquireResponse()
		client := &fasthttp.Client{}
		client.Do(req, resp)

		fmt.Println("Command sent to robot: " + command)
	*/
}

// NewRoboCar constructs object of RoboCar
func NewRoboCar(ip string, port string) (*RoboCar, error) {
	res := &RoboCar{}

	res.CarPort = port
	res.CarIP = ip
	return res, nil
}
