package sender

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	zammadModel "github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	log "github.com/sirupsen/logrus"
	"jaytaylor.com/html2text"
)

type Sender struct {
	vk    *api.VK
	db    *database.DB
	kbrd  *keyboard.Keyboard
	vkCfg config.VKBot
}

func NewSender(vk *api.VK, db *database.DB, kbrd *keyboard.Keyboard, vkCfg config.VKBot) *Sender {
	return &Sender{vk: vk, db: db, kbrd: kbrd, vkCfg: vkCfg}
}

func (s *Sender) Auth(PeerID int) {
	kbrd, err := s.kbrd.GetKeyboard(model.Home, keyboard.Data{})
	if err != nil {
		return
	}

	payload, err := json.Marshal(zammadModel.BotTicket{})
	if err != nil {
		log.Error(err)
		return
	}

	var b = params.NewMessagesSendBuilder()
	b.Payload(string(payload))
	b.Keyboard(string(kbrd))
	b.Message(model.Home.Message())
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(PeerID)
	b.TestMode(true)

	if _, err = s.vk.MessagesSend(b.Params); err != nil {
		log.Error(err)
	}
	return
}

func (s *Sender) Send(whMsg model.WebHookMessage, trigger string) (err error) {
	whMsg.Article.Body, err = html2text.FromString(whMsg.Article.Body, html2text.Options{PrettyTables: true})
	if err != nil {
		log.Error(err)
		return
	}

	var data = model.Data{
		Title: fmt.Sprintf("📄 #%s \"%s\"\n", whMsg.Ticket.Number, whMsg.Ticket.Title),
		WhMsg: whMsg,
	}

	data.Vk, err = s.db.SelectVK(data.WhMsg.Ticket.CustomerID)
	if err != nil {
		log.Error(err)
		return err
	}

	if data.Vk == 0 {
		return
	}

	var b = params.NewMessagesSendBuilder()
	switch trigger {
	case "botNewMessage":
		err = s.botNewMessage(&data)
	case "botChangeGroup":
		err = s.botChangeGroup(&data)
	case "botChangeOwner":
		err = s.botChangeOwner(&data)
	case "botChangeState":
		err = s.botChangeState(&data)
	case "botChangeTitle":
		err = s.botChangeTitle(&data)
	case "botChangePriority":
		err = s.botChangePriority(&data)
	case "botNewTicket":
		err = s.botNewTicket(&data)
	default:
		return
	}

	if err != nil {
		return
	}

	if string(data.Kbrd) != "" {
		b.Keyboard(string(data.Kbrd))
	}

	b.Message(data.Message)
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(data.Vk)
	b.TestMode(true)

	if _, err = s.vk.MessagesSend(b.Params); err != nil {
		log.Error(err)
	}

	return
}

func (s *Sender) botNewMessage(data *model.Data) (err error) {
	data.Message = fmt.Sprintf(
		"%sСообщение от %s: \n\n%s",
		data.Title,
		data.WhMsg.Article.CreatedBy.Displayname,
		data.WhMsg.Article.Body)
	data.Kbrd, err = json.Marshal(model.Keyboard{
		Inline: true,
		Buttons: [][]model.Button{{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.SendMessage.Button(fmt.Sprintf("%d-%d", data.WhMsg.Ticket.ID, data.WhMsg.Ticket.CustomerID))},
				Label:   model.SendMessage.String(),
			},
		}}},
	})
	return
}

func (s *Sender) botNewTicket(data *model.Data) (err error) {
	data.Message = fmt.Sprintf(
		"%sОбращение создано: %s",
		data.Title, data.WhMsg.Article.Body)
	data.Kbrd, err = json.Marshal(model.Keyboard{
		Inline: true,
		Buttons: [][]model.Button{{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.SendMessage.Button(fmt.Sprintf("%d-%d", data.WhMsg.Ticket.ID, data.WhMsg.Ticket.CustomerID))},
				Label:   model.SendMessage.String(),
			},
		}}},
	})
	return
}

func (s *Sender) botChangeGroup(data *model.Data) (err error) {
	if data.WhMsg.Ticket.Group.Name != "" {
		data.Message = fmt.Sprintf("%sИзменена группа: %s.", data.Title, data.WhMsg.Ticket.Group.Name)
	} else {
		data.Message = fmt.Sprintf("%sУдалена группа.", data.Title)
	}
	return
}

func (s *Sender) botChangeOwner(data *model.Data) (err error) {
	if data.WhMsg.Ticket.Owner.Displayname != nil {
		data.Message = fmt.Sprintf(
			"%sИзменен ответственный: %s.",
			data.Title,
			data.WhMsg.Ticket.Owner.Displayname)
	} else if (data.WhMsg.Ticket.Owner.Firstname != "" || data.WhMsg.Ticket.Owner.Lastname != "") &&
		data.WhMsg.Ticket.Owner.Firstname != "-" {
		data.Message = fmt.Sprintf(
			"%sИзменен ответственный: %s %s.",
			data.Title,
			data.WhMsg.Ticket.Owner.Firstname,
			data.WhMsg.Ticket.Owner.Lastname)
	} else {
		data.Message = fmt.Sprintf("%sУдален ответственный.", data.Title)
	}
	return
}

func (s *Sender) botChangeState(data *model.Data) (err error) {
	data.Message = fmt.Sprintf("%sИзменен статус: %s.", data.Title, data.WhMsg.Ticket.State)
	return
}

func (s *Sender) botChangeTitle(data *model.Data) (err error) {
	data.Message = fmt.Sprintf("%sИзменен заголовок.", data.Title)
	return
}

func (s *Sender) botChangePriority(data *model.Data) (err error) {
	data.Message = fmt.Sprintf("%sИзменен приоритет: %s.", data.Title, data.WhMsg.Ticket.Priority.Name)
	return
}
