package lists_transport_http

import (
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

// DeleteList     godoc
// @Summary       Удалить список
// @Description   Удалить существующий в системе список по его ID
// @Tags          lists
// @Param         id path int true "ID удаляемого списка"
// @Success       204 "Успешное удаление списка"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       404 {object} core_http_response.ErrorResponse "List not found"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /lists/{id} [delete]
func (h *ListsHTTPHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
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

	if err := h.listsService.DeleteList(ctx, id); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete list",
		)

		return
	}

	responseHandler.NoContentResponse()
}
