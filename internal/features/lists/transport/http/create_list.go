package lists_transport_http

import (
	"net/http"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type CreateListRequest struct {
	Name         string `json:"name" validate:"required,min=1,max=100" example:"Работа"`
	AuthorUserID int    `json:"author_user_id" validate:"required"     example:"13"`
}

type CreateListResponse ListDTOResponse

// CreateList     godoc
// @Summary       Cоздать список
// @Description   Создать новый список в системе
// @Tags          lists
// @Accept        json
// @Produce       json
// @Param         request body CreateListRequest true "CreateList тело запроса"
// @Success       201 {object} CreateListResponse "Успешно созданный список"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       404 {object} core_http_response.ErrorResponse "Author not found"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /lists [post]
func (h *ListsHTTPHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var req CreateListRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed decode and validate HTTP request",
		)

		return
	}

	domainList := domainFromDTO(req)

	list, err := h.listsService.CreateList(ctx, domainList)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create list",
		)

		return
	}

	response := CreateListResponse(listDTOFromDomain(list))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(req CreateListRequest) domain.List {
	return domain.NewUninitiolizedList(req.Name, req.AuthorUserID)
}
