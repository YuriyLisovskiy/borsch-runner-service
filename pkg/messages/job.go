/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package messages

import "time"

type JobMessage struct {
	ID            string        `json:"id"`
	LangVersion   string        `json:"lang_version"`
	SourceCodeB64 string        `json:"source_code_b64"`
	Timeout       time.Duration `json:"timeout"`
}

type JobResultMessage struct {
	ID       string `json:"id"`
	Data     string `json:"data"`
	ExitCode *int   `json:"exit_code"`
	Error    error  `json:"error"`
}
