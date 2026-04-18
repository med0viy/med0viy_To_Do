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
	Name core_http_types.Nullable[string] `json:"name" swaggertype:"string" example:"Работа"`
}

type PatchListResponse ListDTOResponse

// PatchTask      godoc
// @Summary       Измененить список
// @Description   Измение информации об уже существующем в системе списке
// @Description   ### Логика обновления полей (Three-state logic):
// @Description   1. **Поле не передано**: `name` игнорируется, значение в БД не меняется
// @Description   2. **Явно передано значение**: `"name": "Работа"` - устанавливает новое имя списка в БД
// @Description   3. **Передан null**: `"name": null` - очищает поле в БД (set to NULL)
// @Tags          lists
// @Accept        json
// @Produce       json
// @Param         id path int true "ID изменяемого списка"
// @Param         request body PatchListRequest true "PatchList тело запроса"
// @Success       200 {object} PatchListResponse "Успешно изменнённый список"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       404 {object} core_http_response.ErrorResponse "List not found"
// @Failure       409 {object} core_http_response.ErrorResponse "Сonflict"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /lists/{id} [patch]
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
