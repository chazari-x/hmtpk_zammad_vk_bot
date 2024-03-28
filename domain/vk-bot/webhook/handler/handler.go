package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook/sender"
)

type Handler struct {
	cfg config.WebHook
	s   *sender.Sender
}

func NewHandler(cfg config.WebHook, s *sender.Sender) *Handler {
	return &Handler{cfg: cfg, s: s}
}

// Функция для создания подписи HMAC SHA1
func (wh *Handler) createHmacSignature(data []byte) string {
	h := hmac.New(sha1.New, []byte(wh.cfg.SecretKey))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// WebhookHandler - обработчик вебхука
func (wh *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var message []byte
	var err error

	// Чтение данных из тела запроса
	if message, err = io.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Проверка подписи HMAC SHA1
	receivedSignature := r.Header.Get("X-Hub-Signature")
	expectedSignature := "sha1=" + wh.createHmacSignature(message)
	if receivedSignature != expectedSignature {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Преобразование JSON в структуру
	var Body model.WebHookMessage
	if err = json.Unmarshal(message, &Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Отправка сообщения пользователю
	if err = wh.s.Send(Body, r.Header.Get("X-Zammad-Trigger")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.WriteHeader(http.StatusAccepted)

	return
}
