package lists_transport_http

import (
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetListResponse ListDTOResponse

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
