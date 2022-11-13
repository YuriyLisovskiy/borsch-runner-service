/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package core

import (
	"log"
	"strconv"
)

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

func (jl *JobLogger) OnError(err error) {
	// TODO:
	log.Printf("[JOB ERROR]: %v\n", err)
}

func (jl *JobLogger) OnExit(exitCode int, exitErr error) {
	if exitErr != nil {
		jl.OnError(exitErr)
	}

	jobResult := JobResultMessage{
		ID:   jl.jobId,
		Type: jobResultExit,
		Data: strconv.Itoa(exitCode),
	}
	err := jl.amqpJobService.PublishResult(&jobResult)
	if err != nil {
		jl.OnError(exitErr)
	}
}
