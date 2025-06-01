package responses

import (
	"net/http"
	"user-service/constants"
	errorConstants "user-service/constants/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status   string      `json:"status"`
	Messsage any         `json:"message"`
	Data     interface{} `json:"data"`
	Token    *string     `json:"token,omitempty"`
}

type ParamHTTPResponse struct {
	Code     int
	Err      error
	Messsage *string
	Gin      *gin.Context
	Data     interface{}
	Token    *string
}

func HttpResponse(param ParamHTTPResponse) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:   constants.Success,
			Messsage: http.StatusText(http.StatusOK),
			Data:     param.Data,
			Token:    param.Token,
		})
	}
	message := errorConstants.ErrInternalServerError.Error()
	if param.Messsage != nil {
		message = *param.Messsage
	} else if param.Err != nil {
		if errorConstants.ErrMapping(param.Err) {
			message = param.Err.Error()
		}
	}

	param.Gin.JSON(param.Code, Response{
		Status:   constants.Error,
		Messsage: message,
		Data:     param.Data,
	})

	return
}
