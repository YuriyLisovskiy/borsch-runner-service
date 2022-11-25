package main

import (
	"errors"
	"log"
	"os"

	"github.com/YuriyLisovskiy/borsch-runner-service/internal"
)

func main() {
	server := os.Getenv(internal.EnvRabbitMQServer)
	if server == "" {
		log.Fatalln(errors.New("RabbitMQ server is not set"))
	}

	mq, err := internal.NewRabbitMQService(server)
	if err != nil {
		log.Fatalln(err)
	}

	err = mq.ConsumeJobs()
	if err != nil {
		log.Fatalln(err)
	}
}
