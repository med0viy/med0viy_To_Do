package statistics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetSattisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksComplited             int      `json:"tasks_complited"`
	TasksComplitedRate         *float64 `json:"tasks_complited_rate"`
	TasksAvarageCompletionTime *string  `json:"tasks_avarage_completion_time"`
}

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

func toDTOFromDomain(statistics domain.Statistics) GetSattisticsResponse {
	var avgTime *string
	if statistics.TasksAvarageCompletionTime != nil {
		duration := statistics.TasksAvarageCompletionTime.String()
		avgTime = &duration
	}

	return GetSattisticsResponse{
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
