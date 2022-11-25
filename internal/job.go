/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-runner-service/pkg/docker"
)

type Job struct {
	container *docker.Container
	command   []string
}

func NewJob(image, shell, command, code string, outWriter, errWriter docker.ContainerLogger) *Job {
	return &Job{
		container: &docker.Container{
			Image:  image,
			Stdout: outWriter,
			Stderr: errWriter,
		},
		command: []string{shell, "-c", prepareShellScript(command, code)},
	}
}

func (j *Job) Run() (int, error) {
	if j.container == nil {
		return -1, errors.New("docker container is nil")
	}

	return j.container.Run(j.command...)
}

func prepareShellScript(command, code string) string {
	code = strings.ReplaceAll(code, "\"", "\\\"")
	return strings.ReplaceAll(command, "<code>", fmt.Sprintf("\"%s\"", code))
}
