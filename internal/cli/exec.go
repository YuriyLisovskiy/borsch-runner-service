package cli

import (
	"errors"
	"os"

	"YuriyLisovskiy/borsch-runner-service/internal"
)

func ExecuteApp() error {
	server := os.Getenv(internal.EnvRabbitMQServer)
	if server == "" {
		return errors.New("RabbitMQ server is not set")
	}

	mq, err := internal.NewRabbitMQService(server)
	if err != nil {
		return err
	}

	err = mq.ConsumeJobs()
	if err != nil {
		return err
	}

	return nil
}
