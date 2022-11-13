package main

import (
	"log"
	"os"

	"borsch-runner-service/core"
)

func main() {
	server := os.Getenv("RABBITMQ_SERVER")
	mq, err := core.NewRabbitMQService(server)
	if err != nil {
		log.Fatalln(err)
	}

	err = mq.ConsumeJobs()
	if err != nil {
		log.Fatalln(err)
	}
}
