package users_transport_http

import (
	"net/http"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"                  example:"Burmaldak Burmaldakovich"`
	PhoneNumber *string `json:"phone_nuber" validate:"omitempty,min=10,max=15,startswith=+"  example:"+79998887766"`
}

type CreateUserResponse UserDTOResponse

// CreateUser     godoc
// @Summary       Cоздать пользователя
// @Description   Создать нового пользователя в системе
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         request body CreateUserRequest true "CreateUser тело запроса"
// @Success       201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /users [post]
func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var req CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(req)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create user",
		)

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitiolized(dto.FullName, dto.PhoneNumber)
}
