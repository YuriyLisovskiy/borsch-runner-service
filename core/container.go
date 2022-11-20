/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

const (
	EnvContainerImageTemplate   = "CONTAINER_IMAGE_TEMPLATE"
	EnvContainerShell           = "CONTAINER_SHELL"
	EnvContainerCommandTemplate = "CONTAINER_COMMAND_TEMPLATE"
)

type ContainerLogger interface {
	Log(out string)
}

type DockerContainer struct {
	Image  string
	Stdout ContainerLogger
	Stderr ContainerLogger

	cmd *exec.Cmd
}

func (dc *DockerContainer) Run(args ...string) (int, error) {
	defer func() {
		dc.cmd = nil
	}()

	if dc.cmd != nil {
		return -1, errors.New(fmt.Sprintf("Container %s is already running", dc.Image))
	}

	dc.cmd = exec.Command("docker", append([]string{"run", "--rm", dc.Image}, args...)...)
	stdoutReader, err := dc.cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}

	stdoutScanner, err := newStdScanner(stdoutReader, dc.Stdout)
	if err != nil {
		return -1, err
	}

	stderrReader, err := dc.cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	stderrScanner, err := newStdScanner(stderrReader, dc.Stderr)
	if err != nil {
		return -1, err
	}

	if err := dc.cmd.Start(); err != nil {
		return -1, err
	}

	<-stdoutScanner.doneChan
	<-stderrScanner.doneChan

	err = dc.cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), err
		}

		return -1, err
	}

	return 0, nil
}

type stdScanner struct {
	scanner  *bufio.Scanner
	doneChan <-chan bool
}

func newStdScanner(pipe io.ReadCloser, writer ContainerLogger) (*stdScanner, error) {
	scanner := bufio.NewScanner(pipe)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			writer.Log(scanner.Text())
		}

		done <- true
	}()

	return &stdScanner{
		scanner:  scanner,
		doneChan: done,
	}, nil
}
