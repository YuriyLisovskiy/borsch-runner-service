/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package core

import "log"

type JobLogger struct {
	jobId          string
	amqpJobService AMQPJobService
}

func (jl *JobLogger) Log(output string) {
	notification := JobResultMessage{
		ID:   jl.jobId,
		Type: jobResultLog,
		Data: output,
	}
	err := jl.amqpJobService.PublishResult(&notification)
	if err != nil {
		log.Println(err)
		return
	}
}
