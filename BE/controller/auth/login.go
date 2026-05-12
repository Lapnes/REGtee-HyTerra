package auth

import (
	"errors"
	"log"
	"net/http"

	constant "backend/constant"
	httpAuth "backend/entity/http/auth"
	"backend/util/rest"
	"backend/util/validator"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary Login User
// @Description Generating jwt Session Token
// @Tags Auth
// @Param Payload body auth.Login true "Payload"
// @Success 200 {string} string "Ok"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Failure 423 {string} string "Locked"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/auth/login [POST]
func (ctrl Controller) Login(ctx *gin.Context) {
	req := httpAuth.Login{}
	err := rest.BindJSON(ctx, &req)
	if err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, map[string]string{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}
	if err := validator.Validator.Struct(req); err != nil {
		rest.ResponseError(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := ctrl.svc.Login(req)
	if errors.Is(err, constant.ErrWrongCredentials) {
		rest.ResponseData(ctx, http.StatusBadRequest, map[string]string{
			"password": err.Error(),
		})
		return
	} else if errors.Is(err, constant.ErrUserNotFound) {
		rest.ResponseMessage(ctx, http.StatusNotFound)
		return
	} else if err != nil {
		rest.ResponseMessage(ctx, http.StatusInternalServerError)
		log.Println("err:", err.Error())
		return
	}

	rest.ResponseData(ctx, http.StatusOK, res)
}
