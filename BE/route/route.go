package route

import (
	"os"
	"sync"

	constant "backend/constant"
	v1 "backend/controller"
	"backend/database"
	"backend/docs"
	"backend/util/mqtt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func LoadRouter(router *gin.Engine) (*mqtt.Client, error) {
	swaggerRouter(router)

	var master, slave *gorm.DB
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() { defer wg.Done(); master = database.DBMaster() }()
	go func() { defer wg.Done(); slave = database.DBSlave() }()
	wg.Wait()

	db := &database.DB{Master: master, Slave: slave}
	loader := v1.NewRouterLoader()
	tlsOn := true
	mqttCfg := mqtt.Config{
		BrokerURL:       constant.MQTT_BROKER_URL,
		ClientID:        constant.MQTT_CLIENT_ID,
		Username:        constant.MQTT_USERNAME,
		Password:        constant.MQTT_PASSWORD,
		SubscribeTopics: []string{constant.MQTT_TOPIC_IN},
		QOS:             1,
		UseTLS:          &tlsOn,
	}

	client, err := mqtt.NewClient(mqttCfg, loader.OnMessage)
	if err != nil {
		return nil, err
	}

	loader.Router(router, db, client)

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func swaggerRouter(router *gin.Engine) {
	docs.SwaggerInfo.Title = "HyTerra Backend API"
	docs.SwaggerInfo.Description = "Soil moisture monitoring and irrigation control API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = constant.SWAGGER_HOST
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	if os.Getenv("ENV") != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
