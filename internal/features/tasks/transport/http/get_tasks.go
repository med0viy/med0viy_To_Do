package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	filter, err := getTasksFilterFromRequest(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(
		ctx,
		filter.UserID,
		filter.ListID,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)

		return
	}

	response := GetTasksResponse(TasksDTOFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getTasksFilterFromRequest(r *http.Request) (core_http_request.GetTasksFilter, error) {
	var filter core_http_request.GetTasksFilter

	const (
		userIDQueryParamKey = "user_id"
		listIDQueryParamKey = "list_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return core_http_request.GetTasksFilter{}, fmt.Errorf("get `user_id` query param: %w", err)
	}

	filter.UserID = userID

	listID, err := core_http_request.GetIntQueryParam(r, listIDQueryParamKey)
	if err != nil {
		return core_http_request.GetTasksFilter{}, fmt.Errorf("get `list_id` query param: %w", err)
	}

	filter.ListID = listID

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return core_http_request.GetTasksFilter{}, fmt.Errorf("get `limit` query param: %w", err)
	}

	filter.Limit = limit

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return core_http_request.GetTasksFilter{}, fmt.Errorf("get `offset` query param: %w", err)
	}

	filter.Offset = offset

	return  filter, nil
}
