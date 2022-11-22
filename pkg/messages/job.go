/*
 * Borsch Runner Service
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package messages

type JobMessage struct {
	ID            string `json:"id"`
	LangVersion   string `json:"lang_version"`
	SourceCodeB64 string `json:"source_code_b64"`
}

type JobResultType string

const (
	JobResultLog  JobResultType = "log"
	JobResultExit               = "exit"
)

type JobResultMessage struct {
	ID   string        `json:"id"`
	Type JobResultType `json:"type"`
	Data string        `json:"data"`
}
