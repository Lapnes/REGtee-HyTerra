package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"reflect"

	"github.com/gin-gonic/gin"
)

type BindError struct {
	Message interface{}
	Err     error
}

func BindJSON(ctx *gin.Context, v interface{}) (err error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	// Restore body untuk middleware lain
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// ✅ FIX: v sudah pointer (*UpdateSensorStatus), nggak perlu &v
	err = json.Unmarshal(body, v)
	return err
}

func BindQuery(ctx *gin.Context, v interface{}) (err error) {
	return ctx.BindQuery(v)
}

func BindMultipartFormData(ctx *gin.Context, v interface{}) (err error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		err = errors.New("not a struct")
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}

	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		tagForm := t.Field(i).Tag.Get("form")
		tagFormFile := t.Field(i).Tag.Get("form-file")
		fieldType := val.Field(i).Type()

		if fieldType == reflect.TypeOf([]string{}) {
			val.Field(i).Set(reflect.ValueOf(form.Value[tagForm]))
		} else if fieldType == reflect.TypeOf("") && len(form.Value[tagForm]) > 0 {
			val.Field(i).SetString(form.Value[tagForm][0])
		}

		var file *multipart.FileHeader
		if fieldType == reflect.TypeOf([]*multipart.FileHeader{}) {
			val.Field(i).Set(reflect.ValueOf(form.File[tagFormFile]))
		} else if fieldType == reflect.TypeOf(file) && len(form.File[tagFormFile]) > 0 {
			val.Field(i).Set(reflect.ValueOf(form.File[tagFormFile][0]))
		}
	}

	v = val
	return
}
