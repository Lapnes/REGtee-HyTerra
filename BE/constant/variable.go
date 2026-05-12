package constant

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	DB_HOST            string
	DB_PORT            string
	DB_NAME            string
	DB_USER_MASTER     string
	DB_PASSWORD_MASTER string
	DB_USER_SLAVE      string
	DB_PASSWORD_SLAVE  string
	DB_MASTER          = "Master"
	DB_SLAVE           = "Slave"

	PASSWORD_SALT string
	PRIVATE_KEY   string

	SMTP_HOST     string
	SMTP_PORT     int
	SMTP_AUTH     string
	SMTP_PASSWORD string

	SWAGGER_HOST string

	SERVER_LOCATION *time.Location

	MQTT_BROKER_URL string
	MQTT_CLIENT_ID  string
	MQTT_USERNAME   string
	MQTT_PASSWORD   string
	MQTT_TOPIC_IN   string
	MQTT_TOPIC_OUT  string
	MQTT_KEY_PEM    string
	MQTT_CERT_PEM   string
)

func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Println("[warn] .env file not found, falling back to environment variables")
	}
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_NAME = os.Getenv("DB_NAME")
	DB_USER_MASTER = os.Getenv("DB_USER_MASTER")
	DB_PASSWORD_MASTER = os.Getenv("DB_PASSWORD_MASTER")
	DB_USER_SLAVE = os.Getenv("DB_USER_SLAVE")
	DB_PASSWORD_SLAVE = os.Getenv("DB_PASSWORD_SLAVE")

	PASSWORD_SALT = os.Getenv("PASSWORD_SALT")
	PRIVATE_KEY = os.Getenv("PRIVATE_KEY")

	SMTP_HOST = os.Getenv("SMTP_HOST")
	SMTP_PORT, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTP_AUTH = os.Getenv("SMTP_AUTH")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")

	SWAGGER_HOST = os.Getenv("SWAGGER_HOST")
	SERVER_LOCATION, _ = time.LoadLocation(os.Getenv("SERVER_LOCATION"))

	MQTT_BROKER_URL = os.Getenv("MQTT_BROKER")
	MQTT_CLIENT_ID = os.Getenv("CLIENT_ID")
	MQTT_USERNAME = os.Getenv("MQTT_USER")
	MQTT_PASSWORD = os.Getenv("MQTT_PASS")
	MQTT_TOPIC_IN = os.Getenv("TOPIC_SUBSCRIBE")
	MQTT_TOPIC_OUT = os.Getenv("TOPIC_PUBLISH")
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
