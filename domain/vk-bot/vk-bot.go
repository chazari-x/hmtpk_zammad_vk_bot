package vk_bot

import (
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	longpoll2 "github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/longpoll"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/operation"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/webhook"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/security"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/storage"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	log "github.com/sirupsen/logrus"
)

func Start(cfg config.VKBot, z *zammad.Zammad, s *storage.Storage, db *database.DB, sec *security.Security) (err error) {
	vk := api.NewVK(cfg.Token)

	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		return
	}

	cfg.Name = group[0].Name

	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		return
	}

	kbrd := keyboard.NewKeyboardGetter(vk, z, s, db)
	errCh := make(chan error)

	go func() {
		log.Info("Запуск WebHook'а")
		if err = webhook.NewWebHook(cfg.WebHook, vk, db, kbrd, z, cfg, sec).Start(); err != nil {
			if _, ok := <-errCh; ok {
				errCh <- err
				return
			}
		}

		if _, ok := <-errCh; ok {
			errCh <- fmt.Errorf("WebHook остановлен")
			return
		}
	}()

	go func() {
		log.Info("Запуск LongPoll'а")
		oex := operation.NewOperationExecutor(vk, z, kbrd, s, db, cfg.WebHook.OAuth, sec)
		longpoll2.NewLongPoll(lp, vk, oex).MessageEvent().MessageNew()

		if err = lp.Run(); err != nil {
			if _, ok := <-errCh; ok {
				errCh <- err
				return
			}
		}

		if _, ok := <-errCh; ok {
			errCh <- fmt.Errorf("LongPoll остановлен")
			return
		}
	}()

	return <-errCh
}
