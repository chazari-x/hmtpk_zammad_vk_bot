package cmd

import (
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	vkbot "github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/security"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/storage"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "vk-bot",
		Short: "vk-bot",
		Long:  "vk-bot",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := parseConfig(cmd)
			if err != nil {
				log.Fatalf("Ошибка подключения к конфигурации: %v\n", err)
			}

			log.Trace("Подключение к Zammad...")
			z, err := zammad.NewZammad(cfg.Zammad)
			if err != nil {
				log.Fatalf("Ошибка подключения к Zammad: %v\n", err)
			}
			log.Trace("Подключение к Zammad установлено.")

			s := storage.NewStorage(cfg.VKBot, cmd.Context())

			log.Trace("Подключение в базе данных...")
			db, err := database.NewDB(cfg.DB, cmd.Context())
			if err != nil {
				log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
			}
			log.Trace("Подключение к базе данных установлено.")
			defer func() {
				_ = db.DB.Close()
			}()

			log.Trace("Запуск бота для Zammad...")
			defer log.Trace("Бот для Zammad остановлен.")
			if err = vkbot.Start(cfg.VKBot, z, s, db, security.NewSecurity(cfg.Security)); err != nil {
				log.Fatalf("Ошибка: %v\n", err)
			}
		},
	}
	rootCmd.AddCommand(cmd)
}
