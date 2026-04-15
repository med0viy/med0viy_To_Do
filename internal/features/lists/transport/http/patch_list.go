package lists_transport_http

import (
	"net/http"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
	core_http_types "github.com/med0viy/practika/internal/core/transport/http/types"
)

type PatchListRequest struct {
	Name core_http_types.Nullable[string]
}

type PatchListResponse ListDTOResponse

func (h *ListsHTTPHandler) PatchList(w http.ResponseWriter, r *http.Request) {
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

	var req PatchListRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	listPatch := domain.NewListPatch(
		req.Name.ToDomain(),
	)

	list, err := h.listsService.PatchList(ctx, id, listPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch list",
		)

		return
	}

	response := PatchListResponse(listDTOFromDomain(list))

	responseHandler.JSONResponse(response, http.StatusOK)
}
