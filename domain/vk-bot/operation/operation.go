package operation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/security"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/storage"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	zammadModel "github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	r "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Operation struct {
	vk     *api.VK
	zammad *zammad.Zammad
	kbrd   *keyboard.Keyboard
	store  *storage.Storage
	db     *database.DB
	oauth  config.OAuth
	sec    *security.Security
}

func NewOperationExecutor(
	vk *api.VK,
	z *zammad.Zammad,
	k *keyboard.Keyboard,
	s *storage.Storage,
	db *database.DB,
	oauth config.OAuth,
	sec *security.Security,
) *Operation {
	return &Operation{vk: vk, zammad: z, kbrd: k, store: s, db: db, oauth: oauth, sec: sec}
}

type Data struct {
	page           int
	customer       int
	kbrd           []byte
	marshal        []byte
	get            string
	value          string
	prefix         string
	messageText    string
	messageTextTop string
	link           string
	msg            model.Message
	command        model.Command
	ticket         zammadModel.BotTicket
}

func (o *Operation) ExecuteOperation(msg model.Message) (err error) {
	var data = Data{msg: msg}

	// получить последнюю активность
	if data.get, err = o.store.Get(msg.PeerID, model.Status); err != nil && !errors.Is(err, r.Nil) {
		return
	}

	// получить ticket из message payload или создать пустой тикет
	if msg.MessagePayload != "" {
		log.Trace(msg.MessagePayload)
		if err = json.Unmarshal([]byte(msg.MessagePayload), &data.ticket); err != nil {
			data.ticket = zammadModel.BotTicket{}
		}
	} else {
		data.ticket = zammadModel.BotTicket{}
	}

	// получение zammad user id из базы данных; если customer from ticket != zammad user id, то удалить ticket
	{
		customer := data.ticket.Customer
		if data.customer, err = o.db.SelectZammad(msg.PeerID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return
		}
		data.ticket.Customer = strconv.Itoa(data.customer)
		if data.ticket.Article.Body == "" && (customer != data.ticket.Customer || data.ticket.Customer == "0") {
			data.ticket = zammadModel.BotTicket{Customer: data.ticket.Customer}
		}
	}

	// получить префикс и страницу (в случае если выполняется переход по страницам кнопок)
	{
		if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.ChangeGroup.Key()+"-") {
			data.prefix = model.ChangeGroup.Key()
		} else if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.ChangeType.Key()+"-") {
			data.prefix = model.ChangeType.Key()
		} else if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.ChangeOwner.Key()+"-") {
			data.prefix = model.ChangeOwner.Key()
		} else if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.MyTickets.Key()+"-") {
			data.prefix = model.MyTickets.Key()
		}

		if data.prefix != "" {
			if data.page, err = strconv.Atoi(strings.TrimPrefix(msg.ButtonPayload.Button.Key, data.prefix+"-")); err == nil {
				msg.ButtonPayload.Button = model.MorePayload{Key: data.prefix, Value: data.prefix}
			}
		}
	}

	// получить ticket id и customer (в случае, если нажата кнопка под уведомлением из zammad)
	if strings.HasPrefix(data.msg.ButtonPayload.Button.Key, model.SendMessage.Key()) {
		postfix := strings.Replace(data.msg.ButtonPayload.Button.Key, model.SendMessage.Key(), "", 1)
		if postfix != "" {
			elements := strings.Split(postfix, "-")
			if len(elements) == 2 {
				if elements[1] == data.ticket.Customer {
					if data.ticket.ID, err = strconv.Atoi(elements[0]); err != nil {
						return
					}
					data.msg.ButtonPayload.Button.Key = model.SendMessage.Key()
				}
			}
		}
	}

	log.Tracef(
		"message: %s | customer: %d | get: %s | button value: %s | button key: %s",
		data.msg.Text, data.customer, data.get, data.msg.ButtonPayload.Button.Value, data.msg.ButtonPayload.Button.Key,
	)

	if data.customer == 0 {
		if err = o.authorization(&data); err != nil {
			return
		}
	} else if slices.Contains([]string{
		model.Home.Key(), model.DeleteAuth.Key(), model.CancelSend.Key(), model.Cancel.Key(),
		model.Delete.Key(), model.Send.Key(), model.Authorization.Key(),
	}, data.msg.ButtonPayload.Button.Key) {
		switch data.msg.ButtonPayload.Button.Key {
		case model.DeleteAuth.Key():
			if err = o.deleteAuth(&data); err != nil {
				return
			}
		case model.Cancel.Key():
			if err = o.cancel(&data); err != nil {
				return
			}
		case model.CancelSend.Key():
			if err = o.cancelSend(&data); err != nil {
				return
			}
		case model.Home.Key(), model.Authorization.Key():
			if err = o.home(&data); err != nil {
				return
			}
		case model.Delete.Key():
			if err = o.delete(&data); err != nil {
				return
			}
		case model.Send.Key():
			if err = o.send(&data); err != nil {
				return
			}
		}
	} else if data.ticket.ID == 0 {
		if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangeTitle.Key()) {
			if err = o.changeTitle(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangeGroup.Key()) {
			if err = o.changeGroup(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangeBody.Key()) ||
			data.msg.ButtonPayload.Button.Key == "" {
			if err = o.changeBody(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangeType.Key()) {
			if err = o.changeType(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangePriority.Key()) {
			if err = o.changePriority(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangeDepartment.Key()) {
			if err = o.changeDepartment(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.ChangeOwner.Key()) {
			if err = o.changeOwner(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.MyTickets.Key()) {
			if err = o.myTickets(&data); err != nil {
				return
			}
		}
	} else {
		if slices.Contains([]string{data.msg.ButtonPayload.Button.Key, data.get}, model.SendMessage.Key()) ||
			data.msg.ButtonPayload.Button.Key == "" {
			if err = o.sendMessage(&data); err != nil {
				return
			}
		} else if slices.Contains([]string{data.get}, model.MyTickets.Key()) {
			if err = o.myTickets(&data); err != nil {
				return
			}
		}
	}

	if err = o.store.Set(msg.PeerID, model.Status, data.value); err != nil {
		log.Error(err)
		return
	}

	data.messageText = data.messageTextTop + data.command.Message() + data.messageText
	if data.messageText == "" {
		return
	}

	if data.kbrd, err = o.kbrd.GetKeyboard(data.command, keyboard.Data{
		Page:       data.page,
		Department: data.ticket.Department,
		Group:      data.ticket.Group.ID,
		Customer:   data.customer,
		Link:       data.link,
	}); err != nil {
		log.Error(err)
		return
	}

	if data.marshal, err = json.Marshal(data.ticket); err != nil {
		log.Error(err)
		return
	}

	var b = params.NewMessagesSendBuilder()
	b.Keyboard(string(data.kbrd))
	b.Message(data.messageText)
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(msg.PeerID)
	b.Payload(string(data.marshal))
	b.TestMode(true)
	if _, err = o.vk.MessagesSend(b.Params); err != nil {
		log.Error(err)
	}

	return
}

func (o *Operation) deleteAuth(data *Data) (err error) {
	data.link = o.zammadOAuthLink(data)
	data.command = model.DeleteAuth
	data.customer = 0
	data.ticket.Customer = "0"
	return o.db.DeleteUser(data.msg.PeerID)
}

func (o *Operation) cancel(data *Data) (err error) {
	if data.ticket.Group.Name == "" {
		data.messageTextTop = "♻ Вы отменили создание обращения.\n\n"
		data.ticket = zammadModel.BotTicket{}
		data.command = model.Home
	} else {
		if data.ticket.Owner.Name == "" {
			data.ticket.Department = ""
		}
		data.command = model.CreateTicket
		data.messageText = data.ticket.String()
	}
	return
}

func (o *Operation) cancelSend(data *Data) (err error) {
	data.messageTextTop = "♻ Вы отменили отправку сообщения.\n\n"
	data.ticket.ID = 0
	if data.ticket.Group.Name == "" {
		data.command = model.Home
	} else {
		data.command = model.CreateTicket
	}
	return
}

func (o *Operation) home(data *Data) (err error) {
	data.command = model.Home
	data.ticket = zammadModel.BotTicket{}
	return
}

func (o *Operation) delete(data *Data) (err error) {
	data.ticket = zammadModel.BotTicket{}
	data.command = model.Home
	data.messageTextTop = "♻ Вы удалили своё обращение.\n\n"
	return
}

func (o *Operation) send(data *Data) (err error) {
	if _, err = o.zammad.Ticket.Create(data.ticket); err != nil {
		return err
	}
	data.ticket = zammadModel.BotTicket{}
	data.command = model.Home
	data.messageTextTop = "✅ Вы отправили своё обращение.\n\n"
	return
}

func (o *Operation) changeGroup(data *Data) (err error) {
	if data.page != 0 || data.get == "" {
		data.command = model.ChangeGroup
		data.value = data.command.Key()
		return
	}

	if data.ticket.Owner.Name != "" && data.ticket.Group.Name != data.msg.ButtonPayload.Button.Value {
		data.ticket.Department = ""
		data.ticket.Owner = zammadModel.Owner{}
		data.messageTextTop += "⚠ Ответственный был сброшен из-за изменения группы!\n\n"
	}

	data.command = model.CreateTicket
	data.ticket.Group.Name = data.msg.ButtonPayload.Button.Value
	if data.ticket.Group.ID, err = strconv.Atoi(data.msg.ButtonPayload.Button.Key); err != nil {
		return
	}

	data.messageText = data.ticket.String()
	return
}

func (o *Operation) changeOwner(data *Data) (err error) {
	if data.get == "" {
		data.command = model.ChangeDepartment
		data.value = data.command.Key()
		return
	}
	if data.page != 0 {
		data.command = model.ChangeOwner
		data.value = data.command.Key()
		return
	}

	data.command = model.CreateTicket
	data.ticket.Owner = zammadModel.Owner{
		Name: data.msg.ButtonPayload.Button.Value,
		ID:   data.msg.ButtonPayload.Button.Key,
	}
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) changeDepartment(data *Data) (err error) {
	if data.page != 0 || data.get == "" {
		data.command = model.ChangeDepartment
		data.value = data.command.Key()
		return
	}

	data.command = model.ChangeOwner
	data.value = data.command.Key()
	data.ticket.Department = data.msg.ButtonPayload.Button.Value
	return
}

func (o *Operation) changePriority(data *Data) (err error) {
	if data.get == "" {
		data.command = model.ChangePriority
		data.value = data.command.Key()
		return
	}

	data.command = model.CreateTicket
	data.ticket.Priority = data.msg.ButtonPayload.Button.Value
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) myTickets(data *Data) (err error) {
	if data.page != 0 || data.get == "" {
		data.command = model.MyTickets
		data.value = data.command.Key()
		return
	}

	data.command = model.SendMessage
	data.value = data.command.Key()
	return
}

func (o *Operation) changeType(data *Data) (err error) {
	if data.page != 0 || data.get == "" {
		data.command = model.ChangeType
		data.value = data.command.Key()
		return
	}

	data.command = model.CreateTicket
	data.ticket.Type = zammadModel.Type{
		Key:   data.msg.ButtonPayload.Button.Key,
		Value: data.msg.ButtonPayload.Button.Value,
	}
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) changeTitle(data *Data) (err error) {
	if data.get == "" || data.msg.ButtonPayload.Button.Key != "" {
		data.command = model.ChangeTitle
		data.value = data.command.Key()
		return
	}

	if len(data.msg.Text) > 100 {
		data.command = model.ChangeTitle
		data.messageTextTop = "⚠ Превышено количество символов в заголовке!"
		return
	}

	data.command = model.CreateTicket
	data.ticket.Title = data.msg.Text
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) changeBody(data *Data) (err error) {
	if data.get == "" && data.msg.ButtonPayload.Button.Key == model.ChangeBody.Key() {
		data.command = model.ChangeBody
		data.value = data.command.Key()
		return
	}

	if len(data.msg.Text) > 500 {
		if data.ticket.Article.Body != "" {
			data.command = model.ChangeBody
		} else {
			data.command = model.Home
		}
		data.messageTextTop = "⚠ Превышено количество символов в описании!"
		return
	}

	if data.ticket.Title == "" {
		data.ticket.Title = func() string {
			text := strings.SplitAfter(strings.TrimSpace(data.msg.Text), "\n")[0]
			if len(text) > 100 {
				return text[0:100]
			}
			return text
		}()
	}

	data.command = model.CreateTicket
	data.ticket.Article.Body = data.msg.Text

	if data.ticket.Group.Name == "" {
		data.command = model.ChangeGroup
		data.value = data.command.Key()
		data.messageText = data.ticket.Group.Name
	} else {
		data.messageText = data.ticket.String()
	}

	return
}

func (o *Operation) sendMessage(data *Data) (err error) {
	if data.ticket, err = o.zammad.Ticket.TicketById(strconv.Itoa(data.ticket.ID)); err != nil {
		return
	}
	if data.get == "" && data.msg.ButtonPayload.Button.Key == model.SendMessage.Key() ||
		data.ticket.Article.Body == "" && data.msg.Text == "" {
		data.command = model.SendMessage
		data.value = data.command.Key()
		data.messageTextTop = data.ticket.String()
		return
	}

	ticket := data.ticket
	ticket.Article.Body = data.msg.Text
	if _, err = o.zammad.Ticket.SendToTicket(ticket); err != nil {
		return
	}

	data.messageTextTop = "✅ Вы отправили сообщение.\n\n"
	data.ticket.ID = 0
	data.command = model.Home
	return
}

func (o *Operation) authorization(data *Data) (err error) {
	data.link = o.zammadOAuthLink(data)
	data.command = model.Authorization
	return
}

func (o *Operation) zammadOAuthLink(data *Data) string {
	return (&oauth2.Config{
		ClientID:     o.oauth.ClientID,
		ClientSecret: o.oauth.ClientSecret,
		RedirectURL: fmt.Sprintf(
			"%s?user_sign=%d_%s",
			o.oauth.RedirectURL,
			data.msg.PeerID,
			o.sec.CreateHmacSignature([]byte(strconv.Itoa(data.msg.PeerID))),
		),
		Endpoint: oauth2.Endpoint{
			AuthURL:  o.oauth.AuthURL,
			TokenURL: o.oauth.TokenURL,
		},
	}).AuthCodeURL("")
}
