package cmd

import (
	"errors"
	"os"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Config struct {
	VKBot    config.VKBot
	Zammad   config.Zammad
	Security config.Security
}

func parseConfig(_ *cobra.Command) (cfg Config, err error) {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		ForceQuote:                true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		PadLevelText:              true,
	})

	if logLevel := os.Getenv("BOT_LOG_LEVEL"); logLevel != "" {
		var level log.Level
		if level, err = log.ParseLevel(logLevel); err != nil {
			return
		}
		log.SetLevel(level)
	}

	if cfg.VKBot.Token = os.Getenv("BOT_VK_TOKEN"); cfg.VKBot.Token == "" {
		err = errors.New("BOT_VK_TOKEN is nil")
		return
	}

	if cfg.VKBot.Href = os.Getenv("BOT_VK_API_HREF"); cfg.VKBot.Href == "" {
		err = errors.New("BOT_VK_API_HREF is nil")
		return
	}

	if cfg.VKBot.Chat = os.Getenv("BOT_VK_CHAT_HREF"); cfg.VKBot.Chat == "" {
		err = errors.New("BOT_VK_CHAT_HREF is nil")
		return
	}

	cfg.Security.SecretKey = os.Getenv("BOT_WEBHOOK_SECRET_KEY")

	if cfg.VKBot.WebHook.Port = os.Getenv("BOT_WEBHOOK_PORT"); cfg.VKBot.WebHook.Port == "" {
		err = errors.New("BOT_WEBHOOK_PORT is nil")
		return
	}

	if cfg.Zammad.Token = os.Getenv("BOT_ZAMMAD_TOKEN"); cfg.Zammad.Token == "" {
		err = errors.New("BOT_ZAMMAD_TOKEN is nil")
		return
	}

	if cfg.Zammad.Url = os.Getenv("BOT_ZAMMAD_HREF"); cfg.Zammad.Url == "" {
		err = errors.New("BOT_ZAMMAD_HREF is nil")
		return
	}

	if cfg.VKBot.WebHook.OAuth.ClientID = os.Getenv("BOT_ZAMMAD_OAUTH_CLIENT_ID"); cfg.VKBot.WebHook.OAuth.ClientID == "" {
		err = errors.New("BOT_ZAMMAD_OAUTH_CLIENT_ID is nil")
		return
	}

	if cfg.VKBot.WebHook.OAuth.ClientSecret = os.Getenv("BOT_ZAMMAD_OAUTH_CLIENT_SECRET"); cfg.VKBot.WebHook.OAuth.ClientSecret == "" {
		err = errors.New("BOT_ZAMMAD_OAUTH_CLIENT_SECRET is nil")
		return
	}

	if cfg.VKBot.WebHook.OAuth.RedirectURL = os.Getenv("BOT_ZAMMAD_OAUTH_REDIRECT_URL"); cfg.VKBot.WebHook.OAuth.RedirectURL == "" {
		err = errors.New("BOT_ZAMMAD_OAUTH_REDIRECT_URL is nil")
		return
	}

	if cfg.VKBot.WebHook.OAuth.AuthURL = os.Getenv("BOT_ZAMMAD_OAUTH_AUTH_URL"); cfg.VKBot.WebHook.OAuth.AuthURL == "" {
		err = errors.New("BOT_ZAMMAD_OAUTH_AUTH_URL is nil")
		return
	}

	if cfg.VKBot.WebHook.OAuth.TokenURL = os.Getenv("BOT_ZAMMAD_OAUTH_TOKEN_URL"); cfg.VKBot.WebHook.OAuth.TokenURL == "" {
		err = errors.New("BOT_ZAMMAD_OAUTH_TOKEN_URL is nil")
		return
	}

	return
}
