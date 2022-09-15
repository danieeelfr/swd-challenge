package api

import (
	"net/http"
	"time"

	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/danieeelfr/swd-challenge/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pborman/uuid"

	apimodels "github.com/danieeelfr/swd-challenge/internal/httpserver/api/apimodels"
	taskservice "github.com/danieeelfr/swd-challenge/internal/service/task"
	userservice "github.com/danieeelfr/swd-challenge/internal/service/user"
	"github.com/danieeelfr/swd-challenge/pkg/wait"
)

// TaskHandler holds the task handler implementation
type TaskHandler struct {
	TaskService *taskservice.Task
	UserService *userservice.User
	Routes      []*apimodels.Route
	wait        *wait.Wait
}

// NewTaskHandler return a handler implementation or error
func NewTaskHandler(cfg *config.Config, e *echo.Echo, wg *wait.Wait) (*TaskHandler, error) {
	h := new(TaskHandler)
	h.wait = wg

	g := e.Group("")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	config := middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	}

	g.Use(middleware.JWTWithConfig(config))

	ts, err := taskservice.NewTaskService(cfg)
	if err != nil {
		return nil, err
	}

	h.TaskService = ts

	us, err := userservice.NewUserService(cfg)
	if err != nil {
		return nil, err
	}

	h.UserService = us

	h.setPrivateRoutes(g)

	return h, nil

}

func (h *TaskHandler) setPrivateRoutes(g *echo.Group) {
	g.Add(http.MethodGet, "/task", h.get)
	g.Add(http.MethodPost, "/task", h.save)
}

func (h *TaskHandler) get(ctx echo.Context) error {
	var err error
	f := new(models.TaskFilter)

	f.UserID = ctx.QueryParams().Get("user_id")

	if err = ctx.Bind(&f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload format")
	}

	if err = ctx.Validate(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input payload")
	}

	user := extractUserDataFromContext(ctx)

	if user.Role == "TECHNICIAN" {
		if user.ID != f.UserID {
			return echo.NewHTTPError(http.StatusBadRequest, "technician is only able to see, create or update his own performed tasks.")
		}
	}

	tasks, err := h.TaskService.Find(ctx.Request().Context(), f)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) save(ctx echo.Context) error {

	task := new(models.Task)

	if err := ctx.Bind(&task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload format")
	}

	if err := ctx.Validate(task); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input payload")
	}

	user := extractUserDataFromContext(ctx)

	if user.Role == "TECHNICIAN" {
		if user.ID != task.UserID {
			return echo.NewHTTPError(http.StatusBadRequest, "technician is only able to see, create or update his own performed tasks.")
		}
	}

	task.ID = uuid.New()
	task.PerformedAt = time.Now()

	if err := h.TaskService.Save(ctx.Request().Context(), task, user); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, "saved with success!")
}

func extractUserDataFromContext(ctx echo.Context) *models.User {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["user"].(string)
	name := claims["name"].(string)
	role := claims["role"].(string)
	id := claims["id"].(string)

	return &models.User{
		ID:   id,
		Name: name,
		Role: role,
		User: username,
	}
}
