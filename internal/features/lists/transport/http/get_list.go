package lists_transport_http

import (
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetListResponse ListDTOResponse

// GetTask        godoc
// @Summary       Получить список
// @Description   Получение конкретного списка по его ID
// @Tags          lists
// @Produce       json
// @Param         id path int true "ID получаемого списка"
// @Success       200 {object} GetListResponse "Список успешно найден"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       404 {object} core_http_response.ErrorResponse "List not found"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /lists/{id} [get]
func (h *ListsHTTPHandler) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get ID path value",
		)

		return
	}

	list, err := h.listsService.GetList(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get list",
		)

		return
	}

	response := GetListResponse(listDTOFromDomain(list))

	responseHandler.JSONResponse(response, http.StatusOK)
}
