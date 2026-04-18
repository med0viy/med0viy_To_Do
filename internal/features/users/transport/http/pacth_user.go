package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/med0viy/practika/internal/core/domain"
	core_logger "github.com/med0viy/practika/internal/core/logger"
	core_http_request "github.com/med0viy/practika/internal/core/transport/http/request"
	core_http_response "github.com/med0viy/practika/internal/core/transport/http/response"
	core_http_types "github.com/med0viy/practika/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Full_name    core_http_types.Nullable[string] `json:"full_name"     swaggertype:"string" example:"Иван Иванов"` 
	Phone_number core_http_types.Nullable[string] `json:"phone_number"  swaggertype:"string" example:"+76667778899"`
}

// PatchUser      godoc
// @Summary       Измененить пользователя
// @Description   Измение информации об уже существующем в системе пользователе
// @Description   ### Логика обновления полей (Three-state logic):
// @Description   1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description   2. **Явно передано значение**: `"phone_number": "+76667778899"` - устанавливает новый номер в БД
// @Description   3. **Передан null**: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Tags          users
// @Accept        json
// @Produce       json
// @Param         id path int true "ID изменяемого пользователя"
// @Param         request body PatchUserRequest true "PatchUser тело запроса"
// @Success       200 {object} PatchUserResponse "Успешно изменнённый пользователь"
// @Failure       400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure       404 {object} core_http_response.ErrorResponse "User not found"
// @Failure       409 {object} core_http_response.ErrorResponse "Сonflict"
// @Failure       500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router        /users/{id} [patch]
func (r *PatchUserRequest) Validate() error {
	if r.Full_name.Set {
		if r.Full_name.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.Full_name.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.Phone_number.Set {
		if r.Phone_number.Value != nil {
			phoneNumberLen := len([]rune(*r.Phone_number.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.Phone_number.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must startwith '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.LoggerContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.Full_name.ToDomain(),
		request.Phone_number.ToDomain(),
	)
}
