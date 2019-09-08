package api

import (
	"fmt"

	"github.com/RadiumByte/Robocar-HTTP-Server/cmd/web/app"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// WebServer is providing foreign access to the Robocar
type WebServer struct {
	application app.RobocarServer
}

// ProcessCommand pushes new Command to Application for processing
func (server *WebServer) PushCommand(ctx *fasthttp.RequestCtx) {
	commandStr := ctx.UserValue("command").(string)
	server.application.ProcessCommand(commandStr)
}

// Start initializes Web Server, starts application and begins serving
func (server *WebServer) Start(port string) {
	server.application.Start()

	router := fasthttprouter.New()
	router.PUT("/:command", server.PushCommand)

	fmt.Println("Server is starting on port" + port)
	fmt.Println(fasthttp.ListenAndServe(port, router.Handler))
}

// NewWebServer constructs Web Server
func NewWebServer(application app.RobotServer) (*WebServer, error) {
	res := &WebServer{}
	res.application = application

	return res, nil
}
