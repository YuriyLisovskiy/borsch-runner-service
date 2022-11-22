/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"YuriyLisovskiy/borsch-runner-service/pkg/docker"
	"YuriyLisovskiy/borsch-runner-service/pkg/messages"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPJobService interface {
	ConsumeJobs() error
	PublishResult(jobResult *messages.JobResultMessage) error
}

type RabbitMQJobService struct {
	User     string
	Password string
	Server   string

	connection       *amqp.Connection
	jobChannel       *amqp.Channel
	jobResultChannel *amqp.Channel
	jobQueue         amqp.Queue
	jobResultQueue   amqp.Queue
}

const (
	EnvRabbitMQServer      = "RABBITMQ_SERVER"
	EnvRabbitMQJobQueue    = "RABBITMQ_JOB_QUEUE"
	EnvRabbitMQResultQueue = "RABBITMQ_RESULT_QUEUE"
)

func NewRabbitMQService(server string) (*RabbitMQJobService, error) {
	connection, err := amqp.Dial(server)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	service := &RabbitMQJobService{
		Server:     server,
		connection: connection,
	}
	service.jobChannel, service.jobQueue, err = createQueue(connection, os.Getenv(EnvRabbitMQJobQueue))
	if err != nil {
		connection.Close()
		return nil, err
	}

	service.jobResultChannel, service.jobResultQueue, err = createQueue(connection, os.Getenv(EnvRabbitMQResultQueue))
	if err != nil {
		connection.Close()
		return nil, err
	}

	return service, nil
}

func (mq *RabbitMQJobService) ConsumeJobs() error {
	defer mq.connection.Close()
	defer mq.jobChannel.Close()
	defer mq.jobResultChannel.Close()
	messages, err := mq.jobChannel.Consume(
		mq.jobQueue.Name, // queue
		"",               // consumer
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	var forever chan struct{}
	go func() {
		for d := range messages {
			err = mq.processJob(d.Body)
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			err = d.Ack(false)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func (mq *RabbitMQJobService) PublishResult(jobResult *messages.JobResultMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(jobResult)
	if err != nil {
		return err
	}

	err = mq.jobResultChannel.PublishWithContext(
		ctx,
		"",                     // exchange
		mq.jobResultQueue.Name, // routing key
		false,                  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf("PUBLISHING RESULT: %v\n", jobResult.Data)

	return nil
}

func (mq *RabbitMQJobService) processJob(data []byte) error {
	jobMessage := messages.JobMessage{}
	err := json.Unmarshal(data, &jobMessage)
	if err != nil {
		return err
	}

	log.Printf("PROCESSING JOB: %v\n", jobMessage.ID)

	jobLogger := &JobLogger{
		jobId:          jobMessage.ID,
		amqpJobService: mq,
	}

	sourceCode, err := base64.StdEncoding.DecodeString(jobMessage.SourceCodeB64)
	if err != nil {
		return err
	}

	dockerJob := NewJob(
		strings.ReplaceAll(os.Getenv(docker.EnvContainerImageTemplate), "<language_version>", jobMessage.LangVersion),
		os.Getenv(docker.EnvContainerShell),
		os.Getenv(docker.EnvContainerCommandTemplate),
		string(sourceCode),
		jobLogger,
		jobLogger,
	)
	exitCode, err := dockerJob.Run()
	if err != nil {
		return err
	}

	jobResult := messages.JobResultMessage{
		ID:   jobMessage.ID,
		Type: messages.JobResultExit,
		Data: strconv.Itoa(exitCode),
	}

	return mq.PublishResult(&jobResult)
}

func createQueue(connection *amqp.Connection, name string) (*amqp.Channel, amqp.Queue, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open a channel: %v", err)
	}

	// defer channel.Close()

	queue, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to set QoS: %v", err)
	}

	return channel, queue, nil
}
