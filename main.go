package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"

	mast "mastodon-post/masthandle"
)

type Configuration struct {
	ChanID   string
	Email    string
	Password string
}

var config Configuration

type Plugin struct {
	plugin.MattermostPlugin
	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *Configuration
}

func (p *Plugin) OnActivate() error {

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "mastodon",
		AutoComplete:     true,
		AutoCompleteDesc: "Posts to MX0WVV Mastodon account",
	}); err != nil {
		// ToDo Log Error!
	}

	conf := p.API.GetPluginConfig()
	config.ChanID = conf["channel"].(string)
	config.Email = conf["mastodon-user"].(string)
	config.Password = conf["mastodon-password"].(string)
	p.configuration = &config

	return nil

}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	var responseMsg string
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case "mastodon":
		if args.ChannelId != p.configuration.ChanID {
			responseMsg = "Not allowed in this room!"
		} else {
			responseMsg = "The deed is done!"
			err := mast.ProcessMsg(args.Command, p.configuration.Email, p.configuration.Password)
			if err != nil {
				return &model.CommandResponse{
					ResponseType: model.CommandResponseTypeEphemeral,
					Text:         fmt.Sprintf(err.Error()),
				}, nil
			}
		}
	default:
	}
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf(responseMsg),
	}, nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
