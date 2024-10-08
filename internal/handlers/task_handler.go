package handlers

import (
	"context"
	"encoding/json"
	"github.com/bhupeshpandey/task-manager-nashville/internal/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

const (
	// HealthStatusOK healthy status
	HealthStatusOK = "Ok"
	version        = "v1"
)

type TaskHandler struct {
	serviceName string
	version     string

	grpcClient proto.TaskServiceClient
	websocket  *webSocketHandler
}

func NewTaskHandler(grpcClient proto.TaskServiceClient) *TaskHandler {
	return &TaskHandler{
		grpcClient:  grpcClient,
		serviceName: "nashville-task-service",
		version:     version,
		websocket:   newWebSocketHandler(grpcClient),
	}
}

func (h *TaskHandler) getHealthz(c *gin.Context) {
	// Create and send the response
	healthResponse := &healthzGetResponse{
		ServiceName: h.serviceName,
		Version:     h.version,
		Status:      HealthStatusOK,
	}
	c.JSON(http.StatusOK, healthResponse)
}

func (h *TaskHandler) AddServiceRoutes(wsRouter *gin.RouterGroup, updateHandler func(wsRouter *gin.RouterGroup, method string, path string, handler func(c *gin.Context))) {
	updateHandler(wsRouter, http.MethodGet, "/healthz", h.getHealthz)
	updateHandler(wsRouter, http.MethodPost, "/task", h.createTask)
	updateHandler(wsRouter, http.MethodDelete, "/task/:id", h.deleteTask)
	updateHandler(wsRouter, http.MethodPut, "/task/:id", h.updateTask)
	updateHandler(wsRouter, http.MethodGet, "/task/:id", h.getTask)
	updateHandler(wsRouter, http.MethodGet, "/tasks", h.listTasks)
	wsRouter.GET("/task/ws", h.websocket.handleConnections)
}

func (h *TaskHandler) createTask(c *gin.Context) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ParentID    string `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &proto.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		ParentId:    req.ParentID,
	}
	resp, err := h.grpcClient.CreateTask(context.Background(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *TaskHandler) getTask(c *gin.Context) {
	id := c.Param("id")
	req := proto.GetTaskRequest{
		Id: id,
	}
	res, err := h.grpcClient.GetTask(context.Background(), &req)
	if err != nil {
		s, _ := status.FromError(err)
		c.JSON(http.StatusNotFound, gin.H{"error": s})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *TaskHandler) updateTask(c *gin.Context) {
	var req proto.UpdateTaskRequest
	r := c.Request
	id := c.Param("id")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.Id = id

	res, err := h.grpcClient.UpdateTask(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *TaskHandler) deleteTask(c *gin.Context) {
	id := c.Param("id")
	req := &proto.DeleteTaskRequest{Id: id}
	resp, err := h.grpcClient.DeleteTask(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": s})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TaskHandler) listTasks(c *gin.Context) {
	r := c.Request
	vars := r.URL.Query()
	pageVar := vars.Get("page")
	if pageVar == "" {
		pageVar = "0"
	}
	page, _ := strconv.Atoi(pageVar)
	get := vars.Get("pageSize")
	if get == "" {
		get = "50"
	}
	pageSize, _ := strconv.Atoi(get)
	if pageSize > 50 || pageSize < 0 {
		pageSize = 50
	}
	var startTime, endTime time.Time
	if s := vars.Get("startTime"); s != "" {
		startTime, _ = time.Parse(time.RFC3339, s)

		agoTime := TwoDaysAgoTime()
		if startTime.Before(agoTime) {
			startTime = agoTime
		}
	}
	if e := vars.Get("endTime"); e != "" {
		endTime, _ = time.Parse(time.RFC3339, e)

		if endTime.Before(startTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time cannot be lesser than start time"})
			return
		}
	}

	req := &proto.ListTasksRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	if !startTime.IsZero() {
		req.StartTime = startTime.Format(time.RFC3339)
	}
	if !endTime.IsZero() {
		req.EndTime = endTime.Format(time.RFC3339)
	}

	res, err := h.grpcClient.ListTasks(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": s})
		return
	}
	var resp struct {
		Results interface{} `json:"results"`
	}

	resp.Results = res.Tasks
	c.JSON(http.StatusOK, resp)
}

func TwoDaysAgoTime() time.Time {
	utcTime := time.Now().Add(-48 * time.Hour).UTC()
	return utcTime
}
