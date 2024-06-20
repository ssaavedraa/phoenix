package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	KafkaBroker     string
	ConsumerGroupId string
	Topic           string

	SmtpServer   string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
)

func init() {

	env := os.Getenv("ENVIRONMENT")

	if env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	// Kafka configurations
	KafkaBroker = os.Getenv("KAFKA_BROKER")
	log.Println("kafka broker", KafkaBroker)
	ConsumerGroupId = os.Getenv("CONSUMER_GROUP_ID")
	Topic = os.Getenv("TOPIC")

	// SMTP configurations
	SmtpServer = os.Getenv("SMTP_SERVER")
	SmtpPort, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	SmtpUsername = os.Getenv("SMTP_USERNAME")
	SmtpPassword = os.Getenv("SMTP_PASSWORD")
}
