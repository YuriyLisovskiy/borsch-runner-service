/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package core

import (
	"errors"
	"fmt"
	"strings"
)

type JobInfoHandler interface {
	OnError(err error)
	OnExit(exitCode int, err error)
}

type Job struct {
	container   *DockerContainer
	shell       string
	command     string
	code        string
	infoHandler JobInfoHandler
}

func NewJob(image, shell, command, code string, outWriter, errWriter StdLogger, infoHandler JobInfoHandler) *Job {
	return &Job{
		container: &DockerContainer{
			Image:  image,
			Stdout: outWriter,
			Stderr: errWriter,
		},
		shell:       shell,
		command:     command,
		code:        code,
		infoHandler: infoHandler,
	}
}

func (j *Job) Run() (int, error) {
	if j.container == nil {
		return -1, errors.New("docker container is nil")
	}

	code := strings.ReplaceAll(j.code, "\"", "\\\"")
	shellScript := strings.ReplaceAll(j.command, "<code>", fmt.Sprintf("\"%s\"", code))
	// exitCode, err := j.container.Run(j.shell, "-c", shellScript)
	// if j.infoHandler != nil {
	// 	j.infoHandler.OnExit(exitCode, err)
	// }
	//
	return j.container.Run(j.shell, "-c", shellScript)
}
