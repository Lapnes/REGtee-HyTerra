package readings

import (
	"backend/constant"
	httpReadings "backend/entity/http/readings"
	reading "backend/service/reading"
	"backend/util/jwt/jws"
	"backend/util/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadingsController struct {
	svc reading.Servicer
}

func NewReadingsController(svc reading.Servicer) *ReadingsController {
	return &ReadingsController{svc: svc}
}

func (ctrl *ReadingsController) GetReadings(c *gin.Context) {
	req := httpReadings.GetReadings{}
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
	resp, err := ctrl.svc.GetReadings(req)
	if err != nil {
		rest.ResponseMessage(c, http.StatusInternalServerError)
		return
	}

	rest.ResponseData(c, http.StatusOK, resp)
}
