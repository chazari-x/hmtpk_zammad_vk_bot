package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	cfg config.VKBot
	ctx context.Context
}

func NewStorage(cfg config.VKBot, ctx context.Context) *Storage {
	return &Storage{cfg: cfg, ctx: ctx}
}

type Body struct {
	Response []Response `json:"response"`
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Storage) send(href string) (res []byte, err error) {
	href = fmt.Sprintf("%s%s&v=5.199&access_token=%s", s.cfg.Href, href, s.cfg.Token)
	get, err := http.Get(href)
	if err != nil {
		return
	}

	if get.StatusCode != 200 {
		log.Error(href)
		return res, fmt.Errorf(get.Status)
	}

	return io.ReadAll(get.Body)
}

func (s *Storage) Get(userID int, key string) (value string, err error) {
	all, err := s.send(fmt.Sprintf("/method/storage.get?user_id=%d&key=%s", userID, key))
	if err != nil {
		log.Error(err)
		return
	}

	var resp Body
	if err = json.Unmarshal(all, &resp); err != nil {
		log.Error(err)
		return
	}

	if len(resp.Response) == 0 {
		return value, fmt.Errorf("nil response")
	}

	return resp.Response[0].Value, err
}

func (s *Storage) Set(userID int, key, value string) (err error) {
	href := fmt.Sprintf("/method/storage.set?user_id=%d&key=%s", userID, key)
	if value != "" {
		href = fmt.Sprintf("%s&value=%s", href, value)
	}

	if _, err = s.send(href); err != nil {
		log.Error(err)
		return
	}

	return
}
