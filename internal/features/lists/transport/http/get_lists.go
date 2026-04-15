package lists_transport_http

import (
	"net/http"

	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type GetListsResponse []ListDTOResponse

func (h *ListsHTTPHandler) GetLists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userID, err := core_http_request.GetIntQueryParam(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID query params",
		)

		return
	}

	listsDomains, err := h.listsService.GetLists(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get lists",
		)

		return
	}

	response := GetListsResponse(listsDTOFromDmains(listsDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}
