package main

import (
	"flag"
	"log"
	"telegramBot/clients/telegramClient"
	event_consumer "telegramBot/consumer/event-consumer"
	"telegramBot/events/telegram"
	"telegramBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	token := mustToken()

	tgClient := telegramClient.New(tgBotHost, token)
	eventsProcessor := telegram.New(tgClient, files.New(storagePath))
	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
func mustToken() string {
	token := flag.String("token-bot", "", "token for access telegram bot")
	flag.Parse()
	if *token == "" {
		log.Fatal("Token isn't authorized")
	}
	return *token
}
