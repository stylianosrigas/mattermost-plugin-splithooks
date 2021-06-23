package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (p *Plugin) handleWebhook(body io.Reader) {
	p.API.LogInfo("Received Split.io notification")
	var s *SplitHooksNotification
	if err := json.NewDecoder(body).Decode(&s); err != nil {
		p.postHTTPDebugMessage(err.Error())
		return
	}
	p.API.LogInfo("Message to improve splithooks", "msg=", s.ToJSON())

	environmentNames := strings.Split(p.configuration.EnvironmentNames, ",")

	for _, environmentName := range environmentNames {
		if environmentName == s.EnvironmentName {
			attachment := &model.SlackAttachment{
				Title: fmt.Sprintf("Split.io notification for %s in %s environment", s.Name, s.EnvironmentName),
				Fields: []*model.SlackAttachmentField{
					{Title: "Name", Value: s.Name, Short: true},
					{Title: "Type", Value: s.Type, Short: true},
					{Title: "Change Number", Value: s.ChangeNumber, Short: true},
					{Title: "Time", Value: s.Time, Short: true},
					{Title: "EnvironmentName", Value: s.EnvironmentName, Short: true},
					{Title: "SchemaVersion", Value: s.SchemaVersion, Short: true},
					{Title: "Definition", Value: s.Definition, Short: false},
					{Title: "Link", Value: s.Link, Short: false},
					{Title: "Editor", Value: s.Editor, Short: false},
					{Title: "Description", Value: s.Description, Short: false},
				},
			}

			post := &model.Post{
				ChannelId: p.ChannelID,
				UserId:    p.BotUserID,
			}

			model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
			if _, appErr := p.API.CreatePost(post); appErr != nil {
				p.postHTTPDebugMessage(appErr.Message)
				return
			}
		}
	}
}
