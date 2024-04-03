package webhook

import (
	"net/http"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook/handler"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook/sender"
	log "github.com/sirupsen/logrus"
)

type WebHook struct {
	cfg config.WebHook
	h   *handler.Handler
}

func NewWebHook(cfg config.WebHook, vk *api.VK, db *database.DB, kbrd *keyboard.Keyboard) *WebHook {
	return &WebHook{cfg: cfg, h: handler.NewHandler(cfg, sender.NewSender(vk, db, kbrd))}
}

func (wh *WebHook) Start() error {
	// Регистрируем обработчик вебхука
	http.HandleFunc("/webhook", wh.h.WebhookHandler)

	// Запускаем сервер на порту
	log.Infof("WebHook запущен на порту: %s", wh.cfg.Port)
	return http.ListenAndServe(wh.cfg.Port, nil)
}
