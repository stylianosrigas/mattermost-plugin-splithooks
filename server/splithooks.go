package main

import (
	"encoding/json"
)

type SplitHooksNotification struct {
	Name            string                  `json:"name,omitempty"`
	Type            string                  `json:"type,omitempty"`
	ChangeNumber    int64                   `json:"changeNumber,omitempty"`
	Time            int64                   `json:"time,omitempty"`
	Definition      string                  `json:"definition,omitempty"`
	Description     string                  `json:"description,omitempty"`
	Link            string                  `json:"link,omitempty"`
	EnvironmentName string                  `json:"environmentName,omitempty"`
	Editor          string                  `json:"editor,omitempty"`
	SchemaVersion   int                     `json:"schemaVersion,omitempty"`
	Previous        *SplitHooksNotification `json:"previous,omitempty"`
}

func (o *SplitHooksNotification) ToJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}
