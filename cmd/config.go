package cmd

import (
	"os"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Config struct {
	VKBot    config.VKBot
	Zammad   config.Zammad
	Redis    config.Redis
	DB       config.DataBase
	Security config.Security
}

func parseConfig(_ *cobra.Command) Config {
	var cfg Config

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		ForceQuote:                true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		PadLevelText:              true,
	})
	log.SetLevel(log.TraceLevel)

	cfg.VKBot.Token = os.Getenv("VK_TOKEN")
	cfg.VKBot.Href = os.Getenv("VK_API_HREF")
	cfg.VKBot.Chat = os.Getenv("VK_CHAT_HREF")
	cfg.Zammad.Token = os.Getenv("ZAMMAD_TOKEN")
	cfg.Zammad.Url = os.Getenv("ZAMMAD_HREF")

	cfg.Security.SecretKey = os.Getenv("WEBHOOK_SECRET_KEY")
	cfg.VKBot.WebHook.Port = os.Getenv("WEBHOOK_PORT")
	cfg.VKBot.WebHook.OAuth.ClientID = os.Getenv("ZAMMAD_OAUTH_CLIENT_ID")
	cfg.VKBot.WebHook.OAuth.ClientSecret = os.Getenv("ZAMMAD_OAUTH_CLIENT_SECRET")
	cfg.VKBot.WebHook.OAuth.RedirectURL = os.Getenv("ZAMMAD_OAUTH_REDIRECT_URL")
	cfg.VKBot.WebHook.OAuth.AuthURL = os.Getenv("ZAMMAD_OAUTH_AUTH_URL")
	cfg.VKBot.WebHook.OAuth.TokenURL = os.Getenv("ZAMMAD_OAUTH_TOKEN_URL")

	cfg.DB.Name = os.Getenv("POSTGRES_DB")
	cfg.DB.User = os.Getenv("POSTGRES_USER")
	cfg.DB.Port = os.Getenv("POSTGRES_PORT")
	cfg.DB.Host = os.Getenv("POSTGRES_HOST")
	cfg.DB.Pass = os.Getenv("POSTGRES_PASS")

	return cfg
}
