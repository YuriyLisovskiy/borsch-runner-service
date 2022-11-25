/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

const (
	EnvContainerImageTemplate   = "CONTAINER_IMAGE_TEMPLATE"
	EnvContainerShell           = "CONTAINER_SHELL"
	EnvContainerCommandTemplate = "CONTAINER_COMMAND_TEMPLATE"

	ContainerErrCode = -1
)

var ErrContainerTimedOut = errors.New("timeout error")

type ContainerLogger interface {
	Log(out string)
}

type Container struct {
	Image  string
	Stdout ContainerLogger
	Stderr ContainerLogger

	cmd *exec.Cmd
}

func (dc *Container) Run(timeout time.Duration, args ...string) (int, error) {
	defer func() {
		dc.cmd = nil
	}()

	if dc.cmd != nil {
		return ContainerErrCode, errors.New(fmt.Sprintf("Container %s is already running", dc.Image))
	}

	dc.cmd = exec.Command("docker", append([]string{"run", "--rm", dc.Image}, args...)...)
	stdoutReader, err := dc.cmd.StdoutPipe()
	if err != nil {
		return ContainerErrCode, err
	}

	_, err = newStdScanner(stdoutReader, dc.Stdout)
	if err != nil {
		return ContainerErrCode, err
	}

	stderrReader, err := dc.cmd.StderrPipe()
	if err != nil {
		return ContainerErrCode, err
	}

	_, err = newStdScanner(stderrReader, dc.Stderr)
	if err != nil {
		return ContainerErrCode, err
	}

	err = dc.cmd.Start()
	if err != nil {
		return ContainerErrCode, err
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error)
	go func() {
		errChan <- dc.cmd.Wait()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return exitErr.ExitCode(), err
			}

			return ContainerErrCode, err
		}

		return 0, nil
	case <-timeoutContext.Done():
		err = dc.cmd.Process.Kill()
		if err != nil {
			return ContainerErrCode, errors.New("failed to kill process")
		}

		return int(syscall.SIGKILL), ErrContainerTimedOut
	}
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
