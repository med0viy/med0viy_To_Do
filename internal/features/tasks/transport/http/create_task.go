package tasks_transport_http

import (
	"net/http"
	"time"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
	core_http_types "github.com/med0viy/practika/internal/core/transport/http/types"
)

type CreateTaskRequest struct {
	Title        string                `json:"title" validate:"required,min=1,max=100"         example:"Погулять c собакой"`
	Description  *string               `json:"description" validate:"omitempty,min=1,max=1000" example:"Погулять сегодня вечером c собакой"`
	IsImportant  *bool                 `json:"is_important" validate:"required"                example:"false"`
	IsInMyDay    *bool                 `json:"is_in_my_day" validate:"required"                example:"true"`                
	DueDate      *core_http_types.Date `json:"due_date"                                        example:"2026-04-03"`
	ListID       *int                  `json:"list_id"                                         example:"14"`
	AuthorUserID int                   `json:"author_user_id" validate:"required"              example:"16"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask     godoc
// @Summary       Cоздать задачу
// @Description   Создать новую задачу в системе
// @Tags          tasks
// @Accept        json
// @Produce       json
// @Param         request body CreateTaskRequest true "CreateTask тело запроса"
// @Success       201 {object} CreateTaskResponse "Успешно созданная задача"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       404 {object} core_http_response.ErrorResponse "Author not found"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /tasks [post]
func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var req CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskDomain := domain.NewTaskUninitialized(
		req.Title,
		req.Description,
		*req.IsImportant,
		*req.IsInMyDay,
		(*time.Time)(req.DueDate),
		req.ListID,
		req.AuthorUserID,
	)

	task, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}

	response := CreateTaskResponse(TaskDTOFromDomain(task))

	responseHandler.JSONResponse(response, http.StatusCreated)
}