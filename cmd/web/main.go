package main

import (
	"fmt"

	"github.com/RadiumByte/Robocar-HTTP-Server/cmd/web/api"
	"github.com/RadiumByte/Robocar-HTTP-Server/cmd/web/app"
	"github.com/RadiumByte/Robocar-HTTP-Server/cmd/web/ral"
)

func main() {
	Port := ":8080"

	robot, err := ral.NewRoboCar()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Robot connection failure - program stopped")
		return
	}

	application, err := app.NewApplication(robot)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Application start failure - program stopped")
		return
	}

	server, err := api.NewWebServer(application)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Server start failure - program stopped")
		return
	}
	server.Start(Port)
	fmt.Println("Server started")
}
