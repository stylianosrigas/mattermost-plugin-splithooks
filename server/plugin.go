package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	BotUserKey = "SplitHooksBot"
)

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	BotUserID string
	ChannelID string
}

func (p *Plugin) OnActivate() error {
	configuration := p.getConfiguration()

	if err := p.IsValid(configuration); err != nil {
		return err
	}

	if err := p.ensureBotExists(); err != nil {
		return errors.Wrap(err, "failed to ensure bot user exists")
	}

	team, err := p.API.GetTeamByName(p.configuration.Team)
	if err != nil {
		return err
	}

	channel, err := p.API.GetChannelByName(team.Id, p.configuration.Channel, false)
	if err != nil && err.StatusCode == http.StatusNotFound { //nolint
		channelToCreate := &model.Channel{
			Name:        p.configuration.Channel,
			DisplayName: p.configuration.Channel,
			Type:        model.CHANNEL_OPEN,
			TeamId:      team.Id,
			CreatorId:   p.BotUserID,
		}

		newChannel, errChannel := p.API.CreateChannel(channelToCreate)
		if errChannel != nil {
			return errChannel
		}
		p.ChannelID = newChannel.Id
	} else if err != nil {
		return err
	} else {
		p.ChannelID = channel.Id
	}

	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/webhook":
		token := r.URL.Query().Get("token")
		if token == "" || strings.Compare(token, p.configuration.Token) != 0 {
			errorMessage := "Invalid or missing token"
			p.postHTTPDebugMessage(errorMessage)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		p.handleWebhook(r.Body, p.ChannelID, p.BotUserID)
	default:
		p.postHTTPDebugMessage("Invalid URL path")
		http.NotFound(w, r)
	}
}

func (p *Plugin) IsValid(configuration *configuration) error {
	if configuration.Team == "" {
		return fmt.Errorf("must set a Team")
	}

	if configuration.Channel == "" {
		return fmt.Errorf("must set a Channel")
	}

	if configuration.EnvironmentNames == "" {
		return fmt.Errorf("must set an Environment Name")
	}

	if configuration.Token == "" {
		return fmt.Errorf("must set a Token")
	}

	return nil
}

func (p *Plugin) postHTTPDebugMessage(errorMessage string) {
	p.API.LogDebug("Failed to serve HTTP request", "Error message", errorMessage)
}

func (p *Plugin) ensureBotExists() error {
	// Attempt to find an existing bot
	botUserIDBytes, err := p.API.KVGet(BotUserKey)
	if err != nil {
		return err
	}

	if botUserIDBytes == nil {
		// Create a bot since one doesn't exist
		p.API.LogDebug("Creating bot for SplitHooks plugin")

		bot, err := p.API.CreateBot(&model.Bot{
			Username:    "splithooksbot",
			DisplayName: "SplitHooks Bot",
			Description: "Created by the SplitHooks plugin.",
		})
		if err != nil {
			return err
		}

		// // Give it a profile picture
		// err = p.API.SetProfileImage(bot.UserId, plitHookBotImage)
		// if err != nil {
		// 	p.API.LogError("Failed to set profile image for bot", "err", err)
		// }

		p.API.LogDebug("Bot created for SplitHooks plugin")

		// Save the bot ID
		err = p.API.KVSet(BotUserKey, []byte(bot.UserId))
		if err != nil {
			return err
		}
		p.BotUserID = bot.UserId
	} else {
		p.BotUserID = string(botUserIDBytes)
	}

	return nil
}
