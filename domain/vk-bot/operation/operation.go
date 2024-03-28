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
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/keyboard"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/storage"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	zammadModel "github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	r "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Operation struct {
	vk     *api.VK
	zammad *zammad.Zammad
	kbrd   *keyboard.Keyboard
	store  *storage.Storage
	db     *database.DB
}

func NewOperationExecutor(vk *api.VK, z *zammad.Zammad, k *keyboard.Keyboard, s *storage.Storage, db *database.DB) *Operation {
	return &Operation{vk: vk, zammad: z, kbrd: k, store: s, db: db}
}

var (
	allResponses = []string{
		model.ChangeGroup.Key(),
		model.ChangeType.Key(),
		model.ChangePriority.Key(),
		model.ChangeTitle.Key(),
		model.ChangeBody.Key(),
		model.CreateTicket.Key(),
		model.ChangeOwner.Key(),
		model.ChangeDepartment.Key(),
		model.Password.Key(),
		model.Authorization.Key(),
		model.SendMessage.Key(),
	}

	textResponses = []string{
		model.ChangeTitle.Key(),
		model.ChangeBody.Key(),
		model.Password.Key(),
		model.Authorization.Key(),
		model.SendMessage.Key(),
	}
)

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
	msg            model.Message
	command        model.Command
	ticket         zammadModel.Ticket
}

func (o *Operation) ExecuteOperation(msg model.Message) (err error) {
	var data Data
	data.msg = msg

	if data.get, err = o.store.Get(msg.PeerID, model.Status); err != nil && !errors.Is(err, r.Nil) {
		return
	}

	if msg.MessagePayload != "" {
		log.Trace(msg.MessagePayload)
		if err = json.Unmarshal([]byte(msg.MessagePayload), &data.ticket); err != nil {
			data.ticket = zammadModel.Ticket{}
		}
	} else {
		data.ticket = zammadModel.Ticket{}
	}

	if data.ticket.Customer != "" {
		if data.customer, err = strconv.Atoi(data.ticket.Customer); err != nil {
			return
		}
	} else {
		data.ticket.Customer = strconv.Itoa(data.customer)
	}

	if data.customer <= 0 {
		if data.customer, err = o.db.SelectZammad(msg.PeerID); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return
		}

		data.ticket.Customer = strconv.Itoa(data.customer)
	}

	if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.ChangeGroup.Key()+"-") {
		data.prefix = model.ChangeGroup.Key()
	} else if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.ChangeType.Key()+"-") {
		data.prefix = model.ChangeType.Key()
	} else if strings.HasPrefix(msg.ButtonPayload.Button.Key, model.ChangeOwner.Key()+"-") {
		data.prefix = model.ChangeOwner.Key()
	}

	if data.prefix != "" {
		if data.page, err = strconv.Atoi(strings.TrimPrefix(msg.ButtonPayload.Button.Key, data.prefix+"-")); err == nil {
			msg.ButtonPayload.Button = model.MorePayload{Key: data.prefix, Value: data.prefix}
		}
	}

	log.Tracef(
		"customer: %d | get: %s | button value: %s | button key: %s",
		data.customer,
		data.get,
		msg.ButtonPayload.Button.Value,
		msg.ButtonPayload.Button.Key,
	)

	if data.get != "" && (slices.Contains(allResponses, data.get)) &&
		(msg.ButtonPayload.Button.Key != "" || slices.Contains(textResponses, data.get)) {
		switch msg.ButtonPayload.Button.Key {
		case model.CancelAuth.Key():
			data.command = model.CancelAuth
		case model.DeleteAuth.Key():
			if err = o.deleteAuth(&data); err != nil {
				return
			}
		case model.Cancel.Key():
			if err = o.cancel(&data); err != nil {
				return
			}
		case "":
			switch data.get {
			case model.ChangeTitle.Key():
				if err = o.changeTitle(&data); err != nil {
					return
				}
			case model.ChangeBody.Key():
				if err = o.changeBody(&data); err != nil {
					return
				}
			case model.SendMessage.Key():
				return o.sendMessage(&data)
			case model.Password.Key():
				if err = o.password(&data); err != nil {
					return
				}
			default:
				data.command = isAuthorization(data.customer)
			}
		default:
			if data.page != 0 {
				switch msg.ButtonPayload.Button.Key {
				case model.ChangeType.Key():
					data.command = model.ChangeType
					data.value = data.command.Key()
				case model.ChangeGroup.Key():
					data.command = model.ChangeGroup
					data.value = data.command.Key()
				case model.ChangeOwner.Key():
					data.command = model.ChangeOwner
					data.value = data.command.Key()
				case model.ChangeDepartment.Key():
					data.command = model.ChangeDepartment
					data.value = data.command.Key()
				default:
					data.command = isAuthorization(data.customer)
				}

				break
			}

			switch data.get {
			case model.ChangeGroup.Key():
				if data.ticket.Owner.Name != "" {
					data.ticket.Department = ""
					data.ticket.Owner = zammadModel.Owner{}
					data.messageTextTop += "⚠ Ответственный был сброшен из-за изменения отдела!\n\n"
				}

				data.command = model.CreateTicket
				data.ticket.Group.Name = msg.ButtonPayload.Button.Value
				if data.ticket.Group.ID, err = strconv.Atoi(msg.ButtonPayload.Button.Key); err != nil {
					return
				}

				data.messageText = data.ticket.String()
			case model.ChangeType.Key():
				data.command = model.CreateTicket
				data.ticket.Type = zammadModel.Type{
					Key:   msg.ButtonPayload.Button.Key,
					Value: msg.ButtonPayload.Button.Value,
				}
				data.messageText = data.ticket.String()
			case model.ChangePriority.Key():
				data.command = model.CreateTicket
				data.ticket.Priority = msg.ButtonPayload.Button.Value
				data.messageText = data.ticket.String()
			case model.ChangeOwner.Key():
				data.command = model.CreateTicket
				data.ticket.Owner = zammadModel.Owner{
					Name: msg.ButtonPayload.Button.Value,
					ID:   msg.ButtonPayload.Button.Key,
				}
				data.messageText = data.ticket.String()
			case model.ChangeDepartment.Key():
				data.command = model.ChangeOwner
				data.value = data.command.Key()
				data.ticket.Department = msg.ButtonPayload.Button.Value
			default:
				data.command = isAuthorization(data.customer)
			}
		}
	} else if data.customer > 0 {
		switch msg.ButtonPayload.Button.Key {
		case model.Home.Key():
			data.command = model.Home
		case model.ChangeTitle.Key():
			data.command = model.ChangeTitle
			data.value = data.command.Key()
		case model.ChangeBody.Key():
			data.command = model.ChangeBody
			data.value = data.command.Key()
		case model.ChangeGroup.Key():
			data.command = model.ChangeGroup
			data.value = data.command.Key()
		case model.ChangeType.Key():
			data.command = model.ChangeType
			data.value = data.command.Key()
		case model.ChangePriority.Key():
			data.command = model.ChangePriority
			data.value = data.command.Key()
		case model.ChangeOwner.Key():
			data.command = model.ChangeDepartment
			data.value = data.command.Key()
		case model.Cancel.Key():
			data.command = model.CreateTicket
		case model.Delete.Key():
			data.ticket = zammadModel.Ticket{}
			data.command = model.Home
			data.messageTextTop = "♻ Вы удалили своё обращение.\n\n"
		case model.Send.Key():
			if _, err = o.zammad.Ticket.Create(data.ticket); err != nil {
				return
			}
			data.ticket = zammadModel.Ticket{}
			data.command = model.Home
			data.messageTextTop = "✅ Вы отправили своё обращение.\n\n"
		case model.DeleteAuth.Key():
			data.command = model.DeleteAuth
			data.value = data.command.Key()
			data.customer = 0
			data.ticket.Customer = "0"
			if err = o.db.DeleteUser(msg.PeerID); err != nil {
				return
			}
		case model.SendMessage.Key():
			data.command = model.SendMessage
			data.value = data.command.Key()
			if data.marshal, err = json.Marshal(data.ticket); err != nil {
				return
			}
		case "":
			switch msg.Text {
			case model.Authorization.Key():
				data.command = model.Home
				data.value = data.command.Key()
			case model.Home.Key():
				data.command = model.Home
			default:
				data.command = model.CreateTicket
				if data.ticket.Article.Body == "" {
					data.ticket.Title = fmt.Sprintf("Обращение от %s", time.Now().Format(time.RFC822))
					data.ticket.Article.Body = msg.Text
				}

				if data.ticket.Group.Name == "" {
					data.command = model.ChangeGroup
					data.value = data.command.Key()
					data.messageText = data.ticket.Group.Name
				} else {
					data.messageText = data.ticket.String()
				}
			}
		default:
			data.command = model.Home
		}
	} else {
		err = o.authorization(&data)
		if err != nil {
			log.Error(err)
			return
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
	}); err != nil {
		log.Error(err)
		return
	}

	if data.marshal == nil {
		if data.marshal, err = json.Marshal(data.ticket); err != nil {
			log.Error(err)
			return
		}
	}

	var b = params.NewMessagesSendBuilder()
	b.Keyboard(string(data.kbrd))
	b.Message(data.messageText)
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(msg.PeerID)
	b.Payload(string(data.marshal))
	b.TestMode(true)
	_, err = o.vk.MessagesSend(b.Params)
	if err != nil {
		log.Error(err)
	}

	return
}

func isAuthorization(customer int) model.Command {
	if customer <= 0 {
		return model.Authorization
	}

	return model.Home
}

func (o *Operation) deleteAuth(data *Data) (err error) {
	data.command = model.DeleteAuth
	data.customer = 0
	data.ticket.Customer = "0"
	return o.db.DeleteUser(data.msg.PeerID)
}

func (o *Operation) cancel(data *Data) (err error) {
	data.ticket.Customer = "0"

	if data.ticket.Group.Name == "" {
		data.messageTextTop = "♻ Вы отменили создание обращения.\n\n"
		data.ticket = zammadModel.Ticket{}
		data.command = model.Home
		return
	}

	if data.ticket.Owner.Name == "" {
		data.ticket.Department = ""
	}

	data.command = model.CreateTicket
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) changeTitle(data *Data) (err error) {
	data.command = model.CreateTicket
	data.ticket.Title = data.msg.Text
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) changeBody(data *Data) (err error) {
	data.command = model.CreateTicket
	data.ticket.Article.Body = data.msg.Text
	data.messageText = data.ticket.String()
	return
}

func (o *Operation) sendMessage(data *Data) (err error) {
	if err = o.store.Set(data.msg.PeerID, model.Status, ""); err != nil {
		return
	}

	messageText := "Вы отправили сообщение"

	data.ticket.Article.Body = data.msg.Text

	_, err = o.zammad.Ticket.SendToTicket(data.ticket)
	if err != nil {
		return
	}

	kbrd, err := o.kbrd.GetKeyboard(model.Home, keyboard.Data{})
	if err != nil {
		return
	}

	var b = params.NewMessagesSendBuilder()
	b.Keyboard(string(kbrd))
	b.Message(messageText)
	b.RandomID(int(time.Now().Unix()))
	b.PeerID(data.msg.PeerID)
	b.TestMode(true)
	_, err = o.vk.MessagesSend(b.Params)
	return
}

func (o *Operation) authorization(data *Data) (err error) {
	data.command = model.Password
	data.value = data.command.Key()
	err = o.store.Set(data.msg.PeerID, model.User, data.msg.Text)
	return
}

func (o *Operation) password(data *Data) (err error) {
	var get string
	if get, err = o.store.Get(data.msg.PeerID, model.User); err != nil && !errors.Is(err, r.Nil) {
		return
	}

	me, err := o.zammad.User.Me(get, data.msg.Text)
	if err != nil {
		err = nil
		data.command = model.ErrorAuth
		return
	}

	data.customer = me.ID
	data.ticket.Customer = strconv.Itoa(me.ID)
	data.command = model.Home
	data.messageTextTop = "✅ Вы подтвердили свою личность!\n\n"

	return o.db.InsertUser(data.msg.PeerID, data.customer)
}
