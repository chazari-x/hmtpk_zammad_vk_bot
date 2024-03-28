package sender

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	model2 "github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	vk   *api.VK
	db   *database.DB
	kbrd *keyboard.Keyboard
}

func NewSender(vk *api.VK, db *database.DB, kbrd *keyboard.Keyboard) *Sender {
	return &Sender{vk: vk, db: db, kbrd: kbrd}
}

func (s *Sender) Send(B model.WebHookMessage, trigger string) (err error) {
	if B.Article.Body != "" && B.Article.Internal {
		return
	}

	id, err := s.db.SelectVK(B.Ticket.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		log.Error(err)
		return err
	}

	var title = fmt.Sprintf("- - #%s %s - -\n", B.Ticket.Number, B.Ticket.Title)
	var message string
	var kbrd []byte

	var b = params.NewMessagesSendBuilder()
	switch trigger {
	case "botNewMessage":
		message = fmt.Sprintf("%sСообщение от %s: \"%s\"", title, B.Article.CreatedBy.Displayname, B.Article.Body)

		if kbrd, err = s.kbrd.GetKeyboard(model.SendMessage, keyboard.Data{}); err != nil {
			return
		}
	case "botChangeGroup":
		if B.Ticket.Group.Name != "" {
			message = fmt.Sprintf("%sИзменена группа: %s.", title, B.Ticket.Group.Name)
		} else {
			message = fmt.Sprintf("%sУдалена группа.", title)
		}
	case "botChangeOwner":
		if B.Ticket.Owner.Displayname != nil {
			message = fmt.Sprintf("%sИзменен ответственный: %s.", title, B.Ticket.Owner.Displayname)
		} else if (B.Ticket.Owner.Firstname != "" || B.Ticket.Owner.Lastname != "") && B.Ticket.Owner.Firstname != "-" {
			message = fmt.Sprintf("%sИзменен ответственный: %s %s.", title, B.Ticket.Owner.Firstname, B.Ticket.Owner.Lastname)
		} else {
			message = fmt.Sprintf("%sУдален ответственный.", title)
		}
	case "botChangeStatus":
		message = fmt.Sprintf("%sИзменен статус: %s.", title, B.Ticket.State)
	case "botChangeTitle":
		message = fmt.Sprintf("%sИзменен заголовок.", title)
	case "botChangePriority":
		message = fmt.Sprintf("%sИзменен приоритет: %s.", title, B.Ticket.Priority.Name)
	default:
		return
	}

	if string(kbrd) != "" {
		b.Keyboard(string(kbrd))
	}

	if marshal, err := json.Marshal(model2.Ticket{
		Customer: strconv.Itoa(B.Ticket.CustomerID),
		ID:       B.Ticket.ID,
	}); err != nil {
		log.Error(err)
		return
	} else {
		b.Payload(string(marshal))
	}

	b.Message(message)
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(id)
	b.TestMode(true)

	if _, err = s.vk.MessagesSend(b.Params); err != nil {
		log.Error(err)
		return
	}

	return
}
