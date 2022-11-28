package masthandle

import (
	"context"
	"strings"

	"github.com/mattn/go-mastodon"
)

func ProcessMsg(msgText string, Email string, Password string) error {

	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     "https://mastodon.radio",
		ClientName: "mattermost",
		Scopes:     "read write follow",
		Website:    "https://letstalkradio.org",
	})

	if err != nil {
		return err
	}

	m := mastodon.NewClient(&mastodon.Config{
		Server:       "https://mastodon.radio",
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
	})

	err = m.Authenticate(context.Background(), Email, Password)
	if err != nil {
		return err
	}

	msg := mastodon.Toot{}
	msg.Sensitive = false
	msg.Status = strings.TrimPrefix(msgText, "/mastodon")
	_, perr := m.PostStatus(context.Background(), &msg)
	if perr != nil {
		return err
	}
	return nil
}
