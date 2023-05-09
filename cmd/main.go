package main

import (
	"context"
	"fmt"
	"log"

	eventconsumer "github.com/Phaseant/PasswordManagerBot/internal/consumer/eventConsumer"
	"github.com/Phaseant/PasswordManagerBot/internal/events/telegramEvents"
	"github.com/Phaseant/PasswordManagerBot/internal/repository"
	"github.com/Phaseant/PasswordManagerBot/internal/telegram"
	"github.com/spf13/viper"
	"github.com/xlab/closer"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	HOST      = "api.telegram.org"
	batchSize = 100
)

var db *mongo.Client

func main() {
	closer.Bind(clearDB) //close DB connection
	initConfig()

	tg := telegram.New(HOST, getApiToken()) //telegram client

	db, err := repository.InitMongo(repository.Config{ //mongo client
		Username: viper.GetString("MONGODB_USERNAME"),
		Password: viper.GetString("MONGODB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	repo := repository.New(db, viper.GetString("SECRET_KEY")) //repository to work with DB
	eventsProcessor := telegramEvents.New(tg, repo)           //telegram events processor
	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)
	go func() {
		consumer.Start()
		closer.Close() //graceful shutdown
	}()

	closer.Hold()
}

func getApiToken() string {
	token := viper.GetString("TELEGRAM_API_TOKEN")

	if token == "" {
		log.Fatal("token is not provided")
	}

	return token
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func clearDB() {
	fmt.Println("\nBye bye...")
	db.Disconnect(context.Background())
}
