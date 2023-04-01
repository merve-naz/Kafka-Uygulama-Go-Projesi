package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/config"
	docs "gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/docs"
	"gitlab.bulutbilisimciler.com/bb/source-code/certificate-service/handlers"
)

// Path: Certificate Service
// @Title BB Certificate Generator Service API
// @Description bb.app.certificateservice: microservice for certificate.
// @Version 1.0.12
// @Schemes http https
// @BasePath /api-certificates
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	// server name
	APP_NAME = "bb.app.certificateservice"
	// server description
	APP_DESCRIPTION = "bb.app.certificateservice: microservice for layout."
)

func main() {
	// parse application envs
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("INIT: Cannot get current working directory os.Getwd()")
	}
	config.ReadConfig(dir)

	// log env and port like "bb.app.certificateservice  env: dev, port: 9082"
	env := config.C.App.Env
	port := config.C.App.Port
	log.Printf("INIT: %s env: %s, port: %s", APP_NAME, env, port)

	// in app cache
	inAppCache := NewInAppCacheStore(time.Minute)

	// kafka connections
	krCert := NewKafkaConsumerConnection(
		config.C.Broker.Url,
		config.C.Broker.ConsumerGroup,
		config.C.Broker.Topic,
	)
	kwCertDL := NewKafkaProducerConnection(
		config.C.Broker.Url,
		config.C.Broker.ConsumerGroup,
		config.C.Broker.TopicDL,
	)
	kwCert := NewKafkaProducerConnection(
		config.C.Broker.Url,
		config.C.Broker.ConsumerGroup,
		config.C.Broker.Topic,
	)

	// create application service
	certificatesvc := handlers.NewCertificateService(

		inAppCache,
		krCert,
		kwCertDL,
		kwCert,
	)

	// check env and set gin mode
	gin.SetMode(gin.ReleaseMode)
	if env == "prod" || env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	certificatesvc.InitAssets()

	// init handlers

	certificatesvc.InitRouter(router)
	certificatesvc.InitSubcriber()

	// check env and set swagger
	if !(env == "prod" || env == "production") {
		docs.SwaggerInfo.BasePath = handlers.API_PREFIX
		router.GET(handlers.API_PREFIX+"/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// start server
	log.Println("INIT: HTTP Application " + APP_NAME + " started on port " + port)
	log.Println("INIT: SUBSCRIBER[*] started listening on topic " + config.C.Broker.Topic)

	// fatal if error
	log.Fatal(router.Run(":" + port))
}
