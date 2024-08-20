package server

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-nashville/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/bhupeshpandey/task-manager-nashville/internal/handlers"
	"github.com/bhupeshpandey/task-manager-nashville/internal/proto"
	//"github.com/gorilla/mux"
	//"path/to/internal/handler"
)

func NewHTTPServer(grpcClient proto.TaskServiceClient, serverConfig *models.HttpServer) *http.Server {
	taskHandler := handlers.NewTaskHandler(grpcClient, serverConfig.Logger)

	// create the new Gin engine and setup middleware handler chain
	ge := gin.New()

	gin.SetMode(gin.TestMode)

	ge.UseRawPath = true

	wsRouter := ge.Group(fmt.Sprintf("/service/%s", "v1"))

	// route GET /healthz status getHealthz
	taskHandler.AddServiceRoutes(wsRouter, AddServiceRoutes)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      ge,
	}
	return server
}

func AddServiceRoutes(wsRouter *gin.RouterGroup, method string, path string, handler func(c *gin.Context)) {
	wsRouter.Handle(method, path, handler)
}

func Start(server *http.Server) {
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
