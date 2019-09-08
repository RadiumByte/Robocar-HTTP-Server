package app

import (
	"fmt"
	"strconv"
	"sync"

	"image"
	//"image/color"

	"gocv.io/x/gocv"
)

// bufferEraser cleans input videostream from unnecessary frames
func bufferEraser(source *gocv.VideoCapture, m *sync.Mutex) {
	tmp := gocv.NewMat()
	defer tmp.Close()

	for {
		m.Lock()
		_ = source.Read(&tmp)
		m.Unlock()
	}
}

// RobotServer is an interface for accepting income commands from Web Server
type RobotServer interface {
	ProcessCommand(string)

	ChangeBlocking(bool)
	ChangeManual(bool)

	ChangeCascade(int8)
	Start()
}

// RobotAccessLayer is an interface for RAL usage from Application
type RobotAccessLayer interface {
	SendCommand(string)
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	Robot       RobotAccessLayer
	IsManual    bool
	IsBlocked   bool
	CascadeType int8
}

// ChangeBlocking can block/unblock car movements
func (app *Application) ChangeBlocking(mode bool) {
	app.IsBlocked = mode

	if mode {
		fmt.Println("Car is blocked")
	} else {
		fmt.Println("Car is moving")
	}
}

// ChangeBlocking sets current mode of driving
func (app *Application) ChangeManual(mode bool) {
	app.IsManual = mode

	if mode {
		fmt.Println("Car is on manual control")
	} else {
		fmt.Println("Car is driving automatically")
	}
}

// ChangeCascade changes cascade, assigned to the specific sign
// 0 - stop
// 1 - circle
// 2 - trapeze
func (app *Application) ChangeCascade(cascade int8) {
	app.CascadeType = cascade
	if cascade == 0 {
		fmt.Println("Cascade type changed to Stop Sign")
	} else if cascade == 1 {
		fmt.Println("Cascade type changed to Circle Sign")
	} else if cascade == 2 {
		fmt.Println("Cascade type changed to Trapeze Sign")
	}
}

// ProcessCommand parses command and determines what to do with it
func (app *Application) ProcessCommand(command string) {
	if command == "halt" {
		app.ChangeBlocking(true)

	} else if command == "go" {
		app.ChangeBlocking(false)

	} else if command == "manual" {
		app.ChangeManual(true)

	} else if command == "auto" {
		app.ChangeManual(false)

	} else if command == "stopsign" {
		app.ChangeCascade(0)

	} else if command == "circlesign" {
		app.ChangeCascade(1)

	} else if command == "trapezesign" {
		app.ChangeCascade(2)

	} else {
		firstChar := command[0]
		if firstChar == 's' || firstChar == 'f' || firstChar == 'b' {
			if !app.IsBlocked {
				app.Robot.SendCommand(command)
			}
		}
	}
}

// NewApplication constructs Application
func NewApplication(robot RobotAccessLayer) (*Application, error) {
	res := &Application{}
	res.Robot = robot

	return res, nil
}

func (app *Application) ai() {
	webcam, err := gocv.OpenVideoCapture("rtsp://192.168.1.39:8080/video/h264")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()
	fmt.Println("RTSP videostream claimed...")

	var m sync.Mutex
	go bufferEraser(webcam, &m)
	fmt.Println("Buffer eraser started...")

	window := gocv.NewWindow("Autopilot")
	defer window.Close()

	imgCurrent := gocv.NewMat()
	defer imgCurrent.Close()

	cascadeCircle := gocv.NewCascadeClassifier()
	cascadeCircle.Load("circle.xml")

	cascadeStop := gocv.NewCascadeClassifier()
	cascadeStop.Load("stop.xml")

	cascadeTrapeze := gocv.NewCascadeClassifier()
	cascadeTrapeze.Load("trapeze.xml")

	fmt.Printf("Main loop is starting...")
	for {
		if !app.IsManual {
			if !app.IsBlocked {
				m.Lock()
				ok := webcam.Read(&imgCurrent)
				m.Unlock()

				if !ok {
					fmt.Printf("Error while read RTSP: program aborted...")
					return
				}
				if imgCurrent.Empty() {
					continue
				}

				var target []image.Rectangle

				if app.CascadeType == 0 {
					// Stop cascade
					target = cascadeStop.DetectMultiScale(imgCurrent)

				} else if app.CascadeType == 1 {
					// Circle cascade
					target = cascadeCircle.DetectMultiScale(imgCurrent)

				} else if app.CascadeType == 2 {
					// Trapeze cascade
					target = cascadeTrapeze.DetectMultiScale(imgCurrent)
				}

				var command string

				var centroid image.Point
				centroid.X = target[0].Dx() / 2
				centroid.Y = target[0].Dy() / 2

				//frameCenter := imgCurrent.Cols() / 2
				rightBorder := int(float64(imgCurrent.Cols()) * 0.6)
				leftBorder := int(float64(imgCurrent.Cols()) * 0.4)

				if centroid.X >= leftBorder && centroid.X <= rightBorder {
					// Need to ride forward
					command = "S50A"

				} else if centroid.X < leftBorder {
					// Need to steer left
					var steerValue int
					steerValue = (50 * centroid.X) / leftBorder
					steerValueStr := strconv.Itoa(steerValue)
					command = "S" + steerValueStr + "A"

				} else if centroid.X > rightBorder {
					// Need to steer right
					var steerValue int
					steerValue = 100 - ((50 * centroid.X) / (imgCurrent.Cols() - rightBorder))
					steerValueStr := strconv.Itoa(steerValue)
					command = "S" + steerValueStr + "A"
				}

				app.Robot.SendCommand(command)

				window.IMShow(imgCurrent)
				if window.WaitKey(1) >= 0 {
					break
				}
			}
		}
	}
}

// Start initializes AI process
func (app *Application) Start() {
	go app.ai()
}
