package cmd

import (
	"os"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Config struct {
	VKBot  config.VKBot    `yaml:"vk-bot"`
	Zammad config.Zammad   `yaml:"zammad"`
	Redis  config.Redis    `yaml:"redis"`
	DB     config.DataBase `yaml:"db"`
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
	cfg.VKBot.Api.Href = os.Getenv("VK_API_HREF")
	cfg.Zammad.Token = os.Getenv("ZAMMAD_TOKEN")
	cfg.Zammad.Url = os.Getenv("ZAMMAD_HREF")
	cfg.VKBot.WebHook.SecretKey = os.Getenv("WEBHOOK_SECRET_KEY")
	cfg.VKBot.WebHook.Port = os.Getenv("WEBHOOK_PORT")
	cfg.VKBot.Api.SecretKey = os.Getenv("POSTGRES_DB_SECRET_KEY")
	cfg.DB.Name = os.Getenv("POSTGRES_DB")
	cfg.DB.User = os.Getenv("POSTGRES_USER")
	cfg.DB.Port = os.Getenv("POSTGRES_PORT")
	cfg.DB.Host = os.Getenv("POSTGRES_HOST")
	cfg.DB.Pass = os.Getenv("POSTGRES_PASS")

	return cfg
}
