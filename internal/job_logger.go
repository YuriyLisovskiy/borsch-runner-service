/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package internal

import (
	"log"

	"YuriyLisovskiy/borsch-runner-service/pkg/messages"
)

type JobLogger struct {
	jobId          string
	amqpJobService AMQPJobService
}

func (jl *JobLogger) Log(output string) {
	notification := messages.JobResultMessage{
		ID:   jl.jobId,
		Type: messages.JobResultLog,
		Data: output,
	}
	err := jl.amqpJobService.PublishResult(&notification)
	if err != nil {
		log.Println(err)
		return
	}
}
