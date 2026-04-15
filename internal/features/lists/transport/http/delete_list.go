package lists_transport_http

import (
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

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
