package rest

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Result  interface{}       `json:"result,omitempty"`
	Error   string            `json:"error,omitempty"`
	Message string            `json:"message,omitempty"`
	Detail  map[string]string `json:"detail,omitempty"`
	Status  int               `json:"status,omitempty"`
}

type ResponseResult struct {
	Context   *gin.Context
	RequestID string
}
type ErrorDetails map[string]string

func ResponseError(ctx *gin.Context, status int, detail interface{}, msg ...string) ResponseResult {
	if len(msg) == 0 {
		msg = []string{http.StatusText(status)}
	}
	requestId := requestid.Get(ctx)

	response := Response{
		Error:   requestId,
		Message: msg[0],
	}
	if det, isValid := detail.(validator.ValidationErrors); isValid {
		response.Detail = map[string]string{}
		for _, err := range det {
			response.Detail[strings.ToLower(err.Field())] = err.Tag()
		}
	} else if det, ok := detail.(map[string]string); ok {
		response.Detail = det
	} else if det, ok := detail.(*ErrorDetails); ok {
		response.Detail = *det
	} else if det, ok := detail.(string); ok {
		response.Detail = map[string]string{}
		response.Detail["error"] = det
	}

	ctx.JSON(status, response)
	return ResponseResult{
		Context:   ctx,
		RequestID: requestId,
	}
}

func ResponseSuccess(ctx *gin.Context, status int, payload interface{}) ResponseResult {
	ctx.JSON(status, payload)
	requestId := requestid.Get(ctx)
	return ResponseResult{
		Context:   ctx,
		RequestID: requestId,
	}
}

func ResponseMessage(ctx *gin.Context, status int, msg ...string) ResponseResult {
	if len(msg) > 1 {
		log.Println("response cannot contain more than one message")
		log.Println("proceeding with first message only...")
	}
	if len(msg) == 0 {
		msg = []string{http.StatusText(status)}
	} else if status < 200 || status > 299 {
		log.Println("[GOUTILS-debug]", msg[0])
	}

	response := Response{
		Message: msg[0],
	}
	if status < 200 || status > 299 {
		response.Error = requestid.Get(ctx)
	}

	ctx.JSON(status, response)
	return ResponseResult{ctx, response.Error}
}

func ResponseData(ctx *gin.Context, status int, payload interface{}, msg ...string) ResponseResult {
	requestId := requestid.Get(ctx)
	if len(msg) > 1 {
		log.Println("response cannot contain more than one message")
		log.Println("proceeding with first message only...")
	}
	if len(msg) == 0 {
		msg = []string{http.StatusText(status)}
	}

	response := Response{
		Result:  payload,
		Message: msg[0],
	}

	ctx.JSON(status, response)
	return ResponseResult{ctx, requestId}
}
