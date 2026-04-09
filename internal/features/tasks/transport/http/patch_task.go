package tasks_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
	core_http_types "github.com/med0viy/practika/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string]               `json:"title"`
	Description core_http_types.Nullable[string]               `json:"description"`
	Complited   core_http_types.Nullable[bool]                 `json:"complited"`
	IsImportant core_http_types.Nullable[bool]                 `json:"is_important"`
	IsInMyDay   core_http_types.Nullable[bool]                 `json:"is_in_my_day"`
	DueDate     core_http_types.Nullable[core_http_types.Date] `json:"due_date"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 && titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 && descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Complited.Set {
		if r.Complited.Value == nil {
			return fmt.Errorf("`Complited` can't be NULL")
		}
	}

	if r.IsImportant.Set {
		if r.IsImportant.Value == nil {
			return fmt.Errorf("`IsImportant` can't be NULL")
		}
	}

	if r.IsInMyDay.Set {
		if r.IsInMyDay.Value == nil {
			return fmt.Errorf("`IsInMyDay` can't be NULL")
		}
	}

	return nil
}

func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	var req PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(req)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskId ,taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}

	response := PatchTaskResponse(TaskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	var domainDueDate domain.Nullable[time.Time]
	domainDueDate.Set = request.DueDate.Set

	if request.DueDate.Value != nil {
		t := (time.Time)(*request.DueDate.Value)
		domainDueDate.Value = &t
	}

	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Complited.ToDomain(),
		request.IsImportant.ToDomain(),
		request.IsInMyDay.ToDomain(),
		domainDueDate,
	)
}
