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
	"jaytaylor.com/html2text"
)

type Sender struct {
	vk   *api.VK
	db   *database.DB
	kbrd *keyboard.Keyboard
}

func NewSender(vk *api.VK, db *database.DB, kbrd *keyboard.Keyboard) *Sender {
	return &Sender{vk: vk, db: db, kbrd: kbrd}
}

type Data struct {
	whMsg   model.WebHookMessage
	title   string
	message string
	kbrd    []byte
	vk      int
}

func (s *Sender) Send(whMsg model.WebHookMessage, trigger string) (err error) {
	if whMsg.Article.Body != "" && whMsg.Article.Internal {
		return
	}

	whMsg.Article.Body, err = html2text.FromString(whMsg.Article.Body, html2text.Options{PrettyTables: true})
	if err != nil {
		log.Error(err)
		return
	}

	var data = Data{
		title: fmt.Sprintf("- - #%s %s - -\n", whMsg.Ticket.Number, whMsg.Ticket.Title),
		whMsg: whMsg,
	}

	data.vk, err = s.db.SelectVK(data.whMsg.Ticket.CustomerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		log.Error(err)
		return err
	}

	var b = params.NewMessagesSendBuilder()
	switch trigger {
	case "botNewMessage":
		err = s.botNewMessage(&data)
	case "botChangeGroup":
		err = s.botChangeGroup(&data)
	case "botChangeOwner":
		err = s.botChangeOwner(&data)
	case "botChangeStatus":
		err = s.botChangeStatus(&data)
	case "botChangeTitle":
		err = s.botChangeTitle(&data)
	case "botChangePriority":
		err = s.botChangePriority(&data)
	default:
		return
	}

	if err != nil {
		return
	}

	if string(data.kbrd) != "" {
		b.Keyboard(string(data.kbrd))
	}

	marshal, err := json.Marshal(model2.Ticket{
		Customer: strconv.Itoa(data.whMsg.Ticket.CustomerID),
		ID:       data.whMsg.Ticket.ID,
	})
	if err != nil {
		log.Error(err)
		return
	}

	b.Payload(string(marshal))
	b.Message(data.message)
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(data.vk)
	b.TestMode(true)

	if _, err = s.vk.MessagesSend(b.Params); err != nil {
		log.Error(err)
	}

	return
}

func (s *Sender) botNewMessage(data *Data) (err error) {
	data.message = fmt.Sprintf(
		"%sСообщение от %s: \n\n%s",
		data.title,
		data.whMsg.Article.CreatedBy.Displayname,
		data.whMsg.Article.Body)
	data.kbrd, err = s.kbrd.GetKeyboard(model.SendMessage, keyboard.Data{})
	return
}

func (s *Sender) botChangeGroup(data *Data) (err error) {
	if data.whMsg.Ticket.Group.Name != "" {
		data.message = fmt.Sprintf("%sИзменена группа: %s.", data.title, data.whMsg.Ticket.Group.Name)
	} else {
		data.message = fmt.Sprintf("%sУдалена группа.", data.title)
	}
	return
}

func (s *Sender) botChangeOwner(data *Data) (err error) {
	if data.whMsg.Ticket.Owner.Displayname != nil {
		data.message = fmt.Sprintf(
			"%sИзменен ответственный: %s.",
			data.title,
			data.whMsg.Ticket.Owner.Displayname)
	} else if (data.whMsg.Ticket.Owner.Firstname != "" || data.whMsg.Ticket.Owner.Lastname != "") &&
		data.whMsg.Ticket.Owner.Firstname != "-" {
		data.message = fmt.Sprintf(
			"%sИзменен ответственный: %s %s.",
			data.title,
			data.whMsg.Ticket.Owner.Firstname,
			data.whMsg.Ticket.Owner.Lastname)
	} else {
		data.message = fmt.Sprintf("%sУдален ответственный.", data.title)
	}
	return
}

func (s *Sender) botChangeStatus(data *Data) (err error) {
	data.message = fmt.Sprintf("%sИзменен статус: %s.", data.title, data.whMsg.Ticket.State)
	return
}

func (s *Sender) botChangeTitle(data *Data) (err error) {
	data.message = fmt.Sprintf("%sИзменен заголовок.", data.title)
	return
}

func (s *Sender) botChangePriority(data *Data) (err error) {
	data.message = fmt.Sprintf("%sИзменен приоритет: %s.", data.title, data.whMsg.Ticket.Priority.Name)
	return
}
