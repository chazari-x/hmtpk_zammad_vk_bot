package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook/sender"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/security"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Handler struct {
	vkCfg config.VKBot
	cfg   config.WebHook
	s     *sender.Sender
	z     *zammad.Zammad
	db    *database.DB
	sec   *security.Security
}

func NewHandler(cfg config.WebHook, s *sender.Sender, z *zammad.Zammad, db *database.DB, sec *security.Security, vkCfg config.VKBot) *Handler {
	return &Handler{cfg: cfg, s: s, z: z, db: db, sec: sec, vkCfg: vkCfg}
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
	expectedSignature := "sha1=" + wh.sec.CreateHmacSignature(message)
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

func (wh *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	userSign := r.FormValue("user_sign")
	if userSign == "" {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	elements := strings.Split(userSign, "_")
	if len(elements) != 2 {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	userID := elements[0]
	if userID == "" {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	if id == 0 {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	sign := elements[1]
	if sign == "" {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	if sign != wh.sec.CreateHmacSignature([]byte(userID)) {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	zammadOAuthConfig := &oauth2.Config{
		ClientID:     wh.cfg.OAuth.ClientID,
		ClientSecret: wh.cfg.OAuth.ClientSecret,
		RedirectURL: fmt.Sprintf(
			"%s?user_sign=%s_%s",
			wh.cfg.OAuth.RedirectURL,
			userID,
			wh.sec.CreateHmacSignature([]byte(userID)),
		),
		Endpoint: oauth2.Endpoint{
			AuthURL:  wh.cfg.OAuth.AuthURL,
			TokenURL: wh.cfg.OAuth.TokenURL,
		},
	}

	code := r.FormValue("code")
	token, err := zammadOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	if token == nil {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	if !token.Valid() {
		wh.sendErrorPage(w, r, http.StatusBadRequest)
		return
	}

	me, err := wh.z.User.Me(token.AccessToken)
	if err != nil {
		log.Error(err)
		wh.sendErrorPage(w, r, http.StatusInternalServerError)
		return
	}

	if err = wh.db.InsertUser(id, me.ID); err != nil {
		log.Error(err)
		wh.sendErrorPage(w, r, http.StatusInternalServerError)
		return
	}

	go wh.s.Auth(id)

	t, err := template.ParseFiles("files/success.html")
	if err != nil {
		http.ServeFile(w, r, "files/success.html")
		return
	}
	w.WriteHeader(http.StatusAccepted)
	if err = t.Execute(w, &struct {
		Name string
		URL  string
	}{URL: wh.vkCfg.Chat, Name: wh.vkCfg.Name}); err != nil {
		_, _ = w.Write([]byte(http.StatusText(http.StatusAccepted)))
		return
	}
}

func (wh *Handler) sendErrorPage(w http.ResponseWriter, _ *http.Request, status int) {
	w.WriteHeader(status)
	t, err := template.ParseFiles("files/error.html")
	if err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}
	if err = t.Execute(w, &struct {
		Name string
		URL  string
	}{URL: wh.vkCfg.Chat, Name: wh.vkCfg.Name}); err != nil {
		http.Error(w, http.StatusText(status), status)
		return
	}
}
