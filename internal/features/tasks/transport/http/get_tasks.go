package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks       godoc
// @Summary       Список задач
// @Description   Просмотр списка задач с опциональной пагинацией и/или фильтрацией по ID автора/списка
// @Tags          tasks
// @Produce       json
// @Param         user_id query int false "Фильтрация задач по ID автора"
// @Param         list_id query int false "Фильтрация задач из по ID списка"
// @Param         limit query int false "Размер страницы с задачами"
// @Param         offset query int false "Смещение страницы с задачами"
// @Success       200 {object} GetTasksResponse "Список задач успешно отдан"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /tasks [get]
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

	return filter, nil
}
