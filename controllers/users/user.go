package controllers

import (
	"net/http"
	errorWrap "user-service/common/errors"
	"user-service/common/responses"
	"user-service/domain/dto"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{service: service}
}

func (u *UserController) Login(ctx *gin.Context) {
	request := &dto.LoginRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errorWrap.ErrValidationResponse
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code:     http.StatusUnprocessableEntity,
			Messsage: &errMessage,
			Data:     errResponse,
			Err:      err,
			Gin:      ctx,
		})
		return
	}

	user, err := u.service.GetUser().Login(ctx, request)

	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}
	responses.HttpResponse(responses.ParamHTTPResponse{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

func (u *UserController) Register(ctx *gin.Context) {
	request := &dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errorWrap.ErrValidationResponse
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code:     http.StatusUnprocessableEntity,
			Messsage: &errMessage,
			Data:     errResponse,
			Err:      err,
			Gin:      ctx,
		})
		return
	}

	user, err := u.service.GetUser().Register(ctx, request)

	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}
	responses.HttpResponse(responses.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

func (u *UserController) Update(ctx *gin.Context) {
	request := &dto.UpdateRequest{}
	uuid := ctx.Param("uuid")

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(request)

	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errorWrap.ErrValidationResponse
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code:     http.StatusUnprocessableEntity,
			Messsage: &errMessage,
			Data:     errResponse,
			Err:      err,
			Gin:      ctx,
		})
		return
	}

	user, err := u.service.GetUser().Update(ctx, request, uuid)

	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}
	responses.HttpResponse(responses.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := u.service.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	responses.HttpResponse(responses.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	user, err := u.service.GetUser().GetUserByUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		responses.HttpResponse(responses.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	responses.HttpResponse(responses.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}
