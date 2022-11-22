package cli

import (
	"errors"
	"os"

	"github.com/YuriyLisovskiy/borsch-runner-service/internal/core"
)

func ExecuteApp() error {
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
