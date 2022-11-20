package main

import (
	"errors"
	"log"
	"os"

	"borsch-runner-service/core"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	server := os.Getenv(core.EnvRabbitMQServer)
	if server == "" {
		return errors.New("RabbitMQ server is not set")
	}

	mq, err := core.NewRabbitMQService(server)
	if err != nil {
		return err
	}

	err = mq.ConsumeJobs()
	if err != nil {
		return err
	}

	return nil
}
