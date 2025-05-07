package main

import (
	"flag"
	tgClient "github.com/sol1corejz/read-adviser-bot/internal/clients"
	event_consumer "github.com/sol1corejz/read-adviser-bot/internal/consumer/event-consumer"
	telegram "github.com/sol1corejz/read-adviser-bot/internal/events/telegram"
	"github.com/sol1corejz/read-adviser-bot/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "file-storage"
	batchSize   = 100
)

func main() {

	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Println("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopeed", err)
	}
}

func mustToken() string {

	token := flag.String("tg-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
