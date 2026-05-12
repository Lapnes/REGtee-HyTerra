package main

import (
	"backend/constant"
	"backend/database"
	routers "backend/route"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func main() {
	// Check if seed command is requested
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seedCommand()
		return
	}
	startApp()
}

func seedCommand() {
	constant.LoadEnvVariables()

	// Establish master database connection for seeding
	masterDB := database.DBMaster()
	if masterDB == nil {
		log.Fatal("Failed to connect to master database")
	}

	// Run seeding
	if err := database.SeedDatabase(masterDB); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("✓ Seeding completed successfully!")
}

// CORS adds permissive cross-origin headers for development / API use.
func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

func startApp() {
	constant.LoadEnvVariables()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(CORS)
	router.Use(func(c *gin.Context) {
		reqID := requestid.Get(c)
		log.Printf("[REQ:%s] %s %s", reqID, c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("[REQ:%s] DONE %d", reqID, c.Writer.Status())
	})

	mqttClient, err := routers.LoadRouter(router)
	if err != nil {
		log.Fatalf("startup error: %v", err)
	}

	serverAddr := fmt.Sprintf(":%d", 8089)
	log.Printf("server listening on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("router error: %v", err)
	}
	defer mqttClient.Disconnect(250)
}
