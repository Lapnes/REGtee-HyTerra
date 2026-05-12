package route

import (
	authController "backend/controller/auth"
	mqttController "backend/controller/mqtt"
	readingsController "backend/controller/readings"
	sensorController "backend/controller/sensor"
	"backend/database"
	mqttRepo "backend/repository/mqtt"
	readingsRepo "backend/repository/readings"
	sensorRepo "backend/repository/sensor"
	userRepo "backend/repository/user"
	authService "backend/service/auth"
	mqttlistenService "backend/service/mqttlisten"
	mqttpubService "backend/service/mqttpub"
	readingService "backend/service/reading"
	sensorService "backend/service/sensor"
	"backend/util/mqtt"

	middlewareAuth "backend/middleware/auth"

	"github.com/gin-gonic/gin"
)

type RouterLoader struct {
	AuthMiddleware gin.HandlerFunc
	mqttCtrl       *mqttController.Controller
}

func NewRouterLoader() *RouterLoader {
	return &RouterLoader{}
}

func (l *RouterLoader) OnMessage(topic string, payload []byte) {
	if l.mqttCtrl != nil {
		l.mqttCtrl.OnMessage(topic, payload)
	}
}

func (l *RouterLoader) Router(r *gin.Engine, db *database.DB, client *mqtt.Client) {
	l.AuthMiddleware = middlewareAuth.AuthMiddleware()

	uRepo := userRepo.NewRepository(*db)
	rRepo := readingsRepo.NewRepository(*db)
	sRepo := sensorRepo.NewRepository(*db)
	mRepo := mqttRepo.NewMQTTPublisher(client)

	aSvc := authService.NewService(*uRepo)
	rSvc := readingService.NewService(*rRepo)
	mpSvc := mqttpubService.NewCommandService(mRepo)
	// ✅ FIX: Pass pointer sRepo & rRepo, bukan dereference
	mlSvc := mqttlistenService.NewMQTTService(sRepo, rRepo)
	sSvc := sensorService.NewService(sRepo, mpSvc)

	aCtrl := authController.NewController(aSvc)
	rCtrl := readingsController.NewReadingsController(rSvc)
	sCtrl := sensorController.NewSensorController(sSvc)
	l.mqttCtrl = mqttController.NewController(mlSvc)

	v1 := r.Group("v1")
	l.auth(v1, aCtrl)
	l.sensor(v1, sCtrl)
	l.reading(v1, rCtrl)
}

func (l *RouterLoader) auth(router *gin.RouterGroup, handler *authController.Controller) {
	auth := router.Group("auth")
	auth.POST("/login", handler.Login)
}

func (l *RouterLoader) sensor(router *gin.RouterGroup, handler *sensorController.SensorController) {
	sensor := router.Group("sensor")
	sensor.Use(l.AuthMiddleware)
	sensor.GET("", handler.GetSensor)
	sensor.PUT("/update", handler.UpdateSensorStatus)
}

func (l *RouterLoader) reading(router *gin.RouterGroup, handler *readingsController.ReadingsController) {
	reading := router.Group("reading")
	reading.Use(l.AuthMiddleware)
	reading.GET("", handler.GetReadings)
}
