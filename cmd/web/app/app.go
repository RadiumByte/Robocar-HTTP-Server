package app

// RobocarServer is an interface for accepting income commands from Web Server
type RobocarServer interface {
	ProcessCommand(string)
}

// RobotAccessLayer is an interface for RAL usage from Application
type RobotAccessLayer interface {
	SendCommand(string)
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	Robot   RobotAccessLayer
	isFirst bool
}

// ProcessCommand parses command and determines what to do with it
func (app *Application) ProcessCommand(command string) {
	if app.isFirst {
		app.Robot.SendCommand("CONNSERV")
		app.isFirst = false
	}
	app.Robot.SendCommand(command)
}

// NewApplication constructs Application
func NewApplication(robot RobotAccessLayer) (*Application, error) {
	res := &Application{}
	res.Robot = robot
	res.isFirst = true

	return res, nil
}
