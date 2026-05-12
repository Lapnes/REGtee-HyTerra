package sensor

import (
	"backend/constant"
	httpSensor "backend/entity/http/sensor"
	sensor "backend/service/sensor"
	"backend/util/jwt/jws"
	"backend/util/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SensorController struct {
	svc sensor.Servicer
}

func NewSensorController(svc sensor.Servicer) *SensorController {
	return &SensorController{svc: svc}
}

func (ctrl *SensorController) GetSensor(c *gin.Context) {
	req := httpSensor.GetSensor{}
	if err := rest.BindQuery(c, &req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, map[string]interface{}{
			"query": constant.ErrInvalid.Error(),
		})
		return
	}

	claim, err := jws.ExtractClient(c.GetHeader("Authorization"))
	if err != nil || claim.ID < 0 {
		rest.ResponseError(c, http.StatusUnauthorized, map[string]interface{}{
			"authorization": constant.ErrInvalid.Error(),
		})
		return
	}

	req.Pgn.Paginate()
	resp, err := ctrl.svc.GetSensor(req)
	if err != nil {
		rest.ResponseMessage(c, http.StatusInternalServerError)
		return
	}

	rest.ResponseData(c, http.StatusOK, resp)
}

func (ctrl *SensorController) UpdateSensorStatus(c *gin.Context) {
	req := httpSensor.UpdateSensorStatus{}
	if err := rest.BindJSON(c, &req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, map[string]interface{}{
			"body": constant.ErrInvalid.Error(),
		})
		return
	}

	claim, err := jws.ExtractClient(c.GetHeader("Authorization"))
	if err != nil || claim.ID < 0 {
		rest.ResponseError(c, http.StatusUnauthorized, map[string]interface{}{
			"authorization": constant.ErrInvalid.Error(),
		})
		return
	}

	if err := ctrl.svc.UpdateSensorStatus(req); err != nil {
		rest.ResponseMessage(c, http.StatusInternalServerError)
		return
	}

	rest.ResponseSuccess(c, http.StatusAccepted, "Success")
}
