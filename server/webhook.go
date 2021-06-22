package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	splitHooksIconURL  = "https://d25hn4jiqx5f7l.cloudfront.net/companies/logos/original/split-io_1532363931.png"
	splitHooksUsername = "SplitHooks Bot"
)

func (p *Plugin) handleWebhook(body io.Reader, channelID, userID string) {
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
				Title: fmt.Sprintf("This is a Split.io notification for %s", s.Name),
				Fields: []*model.SlackAttachmentField{
					{Title: "Name", Value: s.Name, Short: false},
					{Title: "Type", Value: s.Type, Short: false},
					{Title: "Change Number", Value: s.ChangeNumber, Short: false},
					{Title: "Time", Value: s.Time, Short: false},
					{Title: "Definition", Value: s.Definition, Short: false},
					{Title: "Description", Value: s.Description, Short: false},
					{Title: "Link", Value: s.Link, Short: false},
					{Title: "EnvironmentName", Value: s.EnvironmentName, Short: false},
					{Title: "Editor", Value: s.Editor, Short: false},
					{Title: "SchemaVersion", Value: s.SchemaVersion, Short: false},
				},
			}

			post := &model.Post{
				ChannelId: channelID,
				UserId:    userID,
				Props: map[string]interface{}{
					"from_webhook":      "true",
					"override_username": splitHooksUsername,
					"override_icon_url": splitHooksIconURL,
				},
			}

			model.ParseSlackAttachment(post, []*model.SlackAttachment{attachment})
			if _, appErr := p.API.CreatePost(post); appErr != nil {
				p.postHTTPDebugMessage(appErr.Message)
				return
			}
		}
	}
}
