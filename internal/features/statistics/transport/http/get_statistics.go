package statistics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"                 example:"30"`
	TasksComplited             int      `json:"tasks_complited"               example:"15"`
	TasksComplitedRate         *float64 `json:"tasks_complited_rate"          example:"50"`
	TasksAvarageCompletionTime *string  `json:"tasks_avarage_completion_time" example:"36m20s"`
}

// GetStatistics  godoc
// @Summary       Получение статистики
// @Description   Получение статистики по задачам с опциональной фильтрацией по user_id и/или временному промежутку
// @Tags          statistics
// @Produce       json
// @Param         user_id query int false "Фильтрация статистики по конкретному пользователю"
// @Param         list_id query int false "Фильтрация статистики по конкретному списку"
// @Param         from query string false "Начало промежутка рассмотрения статистики (включительно), формат: YYYY-MM-DD"
// @Param         to   query string false "Конец промежутка рассмотрения статистики (не включительно), формат: YYYY-MM-DD"
// @Success       200 {object} GetStatisticsResponse "Успешное получение статистики"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	filter, err := getStatisticsFilterFromRequest(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get query params",
		)

		return
	}

	statistics, err := h.statisticsService.GetStatistics(
		ctx,
		filter.UserID,
		filter.ListID,
		filter.From,
		filter.To,
	)

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)

		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAvarageCompletionTime != nil {
		duration := statistics.TasksAvarageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksComplited:             statistics.TasksComplited,
		TasksComplitedRate:         statistics.TasksComplitedRate,
		TasksAvarageCompletionTime: avgTime,
	}
}

func getStatisticsFilterFromRequest(r *http.Request) (core_http_request.GetStatisticsFilter, error) {
	var filter core_http_request.GetStatisticsFilter

	const (
		userIDQueryParamKey = "user_id"
		listIDQueryParamKey = "list_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return core_http_request.GetStatisticsFilter{}, fmt.Errorf(
			"get 'user_id' query param: %w",
			err,
		)
	}

	filter.UserID = userID

	listID, err := core_http_request.GetIntQueryParam(r, listIDQueryParamKey)
	if err != nil {
		return core_http_request.GetStatisticsFilter{}, fmt.Errorf(
			"get `list_id` query param: %w",
			err,
		)
	}

	filter.ListID = listID

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return core_http_request.GetStatisticsFilter{}, fmt.Errorf(
			"get 'from' query param: %w",
			err,
		)
	}

	filter.From = from

	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return core_http_request.GetStatisticsFilter{}, fmt.Errorf(
			"get 'to' query param: %w",
			err,
		)
	}

	filter.To = to

	return filter, nil
}
