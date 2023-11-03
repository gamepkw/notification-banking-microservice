package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gopkg.in/Shopify/sarama.v1"

	_notificationService "github.com/gamepkw/notification-banking-microservice/internal/services"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// logger.Info("start program...")

	config := sarama.NewConfig()
	config.ClientID = "my-kafka-client"
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true

	brokers := []string{viper.GetString("kafka.broker_address")}
	kafkaClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaClient.Close()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001", "http://localhost:3000", "http://localhost:8090"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	nu := _notificationService.NewNotificationService(timeoutContext, kafkaClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go nu.SendTransactionNotification(ctx)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint

}
