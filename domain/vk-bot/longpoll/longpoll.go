package longpoll

import (
	"context"
	"encoding/json"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/operation"
	log "github.com/sirupsen/logrus"
)

type LongPoll struct {
	vk *api.VK
	lp *longpoll.LongPoll
	o  *operation.Operation
}

func NewLongPoll(lp *longpoll.LongPoll, vk *api.VK, o *operation.Operation) *LongPoll {
	return &LongPoll{lp: lp, vk: vk, o: o}
}

func (l *LongPoll) MessageNew() *LongPoll {
	l.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		var Payload model.Payload

		if obj.Message.Payload != "" {
			if err := json.Unmarshal([]byte(obj.Message.Payload), &Payload); l.error(obj.Message.PeerID, err) && obj.Message.Text == "" {
				log.Error(err)
				return
			}
		}

		var P = model.Message{
			ButtonPayload: model.ButtonPayload{Button: model.MorePayload{
				Key:   Payload.Button,
				Value: Payload.Button,
			}},
			PeerID: obj.Message.PeerID,
			Text:   obj.Message.Text,
			ID:     obj.Message.ID,
		}

		{
			var b = params.NewMessagesGetHistoryBuilder()
			b.PeerID(obj.Message.PeerID).TestMode(true)
			if history, err := l.vk.MessagesGetHistory(b.Params); l.error(obj.Message.PeerID, err) {
				log.Error(err)
				return
			} else if len(history.Items) > 0 {
				for _, item := range history.Items[0:func() int {
					if history.Count >= 20 {
						return 20
					}
					return history.Count
				}()] {
					if item.Payload != "" {
						P.MessagePayload = item.Payload
						break
					}
				}
			}
		}

		l.error(obj.Message.PeerID, l.o.ExecuteOperation(P))
	})

	return l
}

func (l *LongPoll) MessageEvent() *LongPoll {
	l.lp.MessageEvent(func(ctx context.Context, obj events.MessageEventObject) {
		var more model.MorePayload
		var Payload model.Payload
		var P = model.Message{
			PeerID: obj.PeerID,
			ID:     obj.ConversationMessageID,
		}

		if string(obj.Payload) != "" {
			if err := json.Unmarshal(obj.Payload, &Payload); l.error(obj.PeerID, err) {
				log.Error(err)
				goto eventSuccess
			}
		}

		{
			var b = params.NewMessagesGetHistoryBuilder()
			b.PeerID(obj.PeerID).TestMode(true)
			if history, err := l.vk.MessagesGetHistory(b.Params); l.error(obj.PeerID, err) {
				log.Error(err)
				goto eventSuccess
			} else if len(history.Items) > 0 {
				for _, item := range history.Items[0:func() int {
					if history.Count >= 20 {
						return 20
					}
					return history.Count
				}()] {
					if item.Payload != "" {
						P.MessagePayload = item.Payload
						break
					}
				}
			}
		}
		if err := json.Unmarshal([]byte(Payload.Button), &more); l.error(obj.PeerID, err) {
			log.Error(err)
			goto eventSuccess
		}

		P.ButtonPayload = model.ButtonPayload{Button: more}

		l.error(obj.PeerID, l.o.ExecuteOperation(P))

	eventSuccess:
		{
			var b = params.NewMessagesSendMessageEventAnswerBuilder()
			b.PeerID(obj.PeerID).EventID(obj.EventID).UserID(obj.UserID).TestMode(true)
			if _, err := l.vk.MessagesSendMessageEventAnswer(b.Params); l.error(obj.PeerID, err) {
				log.Error(err)
			}
		}
	})

	return l
}

func (l *LongPoll) error(peerID int, err error) bool {
	if err == nil {
		return false
	}

	var b = params.NewMessagesSendBuilder()
	b.Message("❗ Произошла ошибка. Повторите попытку.").RandomID(int(time.Now().Unix())).PeerID(peerID)
	b.TestMode(true)

	if _, err = l.vk.MessagesSend(b.Params); err != nil {
		log.Error(err)
	}

	return true
}
