package webhook

import (
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook/handler"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook/sender"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/security"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	log "github.com/sirupsen/logrus"
)

type WebHook struct {
	cfg config.WebHook
	h   *handler.Handler
}

func NewWebHook(cfg config.WebHook, vk *api.VK, db *database.DB, kbrd *keyboard.Keyboard, z *zammad.Zammad, vkCfg config.VKBot, sec *security.Security) *WebHook {
	return &WebHook{cfg: cfg, h: handler.NewHandler(cfg, sender.NewSender(vk, db, kbrd, vkCfg), z, db, sec)}
}

func (wh *WebHook) Start() error {
	// Регистрируем обработчик вебхука
	http.HandleFunc("/zammad/webhook", wh.h.WebhookHandler)
	http.HandleFunc("/zammad/auth", wh.h.AuthHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "files/favicon.ico")
	})

	// Запускаем сервер на порту
	log.Infof("WebHook запущен на порту: %s", wh.cfg.Port)
	return http.ListenAndServe(wh.cfg.Port, nil)
}
