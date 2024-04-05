package ticket

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	"github.com/chazari-x/zammad-go"
	log "github.com/sirupsen/logrus"
)

type Ticket struct {
	client *zammad.Client
}

func NewTicketController(client *zammad.Client) *Ticket {
	return &Ticket{client: client}
}

func (t *Ticket) TicketById(id string) (Ticket model.BotTicket, err error) {
	atoi, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
		return
	}

	data, err := t.client.TicketShow(atoi)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var ticket model.Ticket
	if err = json.Unmarshal(bytes, &ticket); err != nil {
		log.Error(err)
		return
	}

	data, err = t.client.TicketPriorityShow(ticket.PriorityID)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err = json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var priority model.Priority
	if err = json.Unmarshal(bytes, &priority); err != nil {
		log.Error(err)
		return
	}

	data, err = t.client.GroupShow(ticket.GroupID)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err = json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var group model.Group
	if err = json.Unmarshal(bytes, &group); err != nil {
		log.Error(err)
		return
	}

	data, err = t.client.UserShow(ticket.OwnerID)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err = json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var owner model.User
	if err = json.Unmarshal(bytes, &owner); err != nil {
		log.Error(err)
		return
	}

	Ticket = model.BotTicket{
		ID:         ticket.ID,
		Title:      ticket.Title,
		Group:      group,
		Customer:   strconv.Itoa(ticket.CustomerID),
		Department: owner.Department,
		Priority:   priority.Name,
		Owner: model.Owner{
			Name: func() string {
				if owner.DisplayName != "" && owner.DisplayName != "-" {
					return owner.DisplayName
				}

				return fmt.Sprintf("%s %s", owner.Firstname, owner.Lastname)
			}(),
			ID: strconv.Itoa(owner.ID),
		},
	}

	return
}

func (t *Ticket) TicketsByCustomer(customer int) (TicketsByCustomer model.TicketsByCustomer, err error) {
	data, err := t.client.TicketListByCustomer(customer)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(bytes, &TicketsByCustomer); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) PriorityList() (Priorities []model.Priority, err error) {
	data, err := t.client.TicketPriorityList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(bytes, &Priorities); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) Create(ticket model.BotTicket) (TicketCreate model.Ticket, err error) {
	ticketInterface, err := ticket.Interface()
	if err != nil {
		log.Error(err)
		return
	}

	data, err := t.client.TicketCreate(ticketInterface)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(bytes, &TicketCreate); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) SendToTicket(ticket model.BotTicket) (Article model.TicketArticle, err error) {
	var article = model.TicketArticleCreate{
		Body:     ticket.Article.Body,
		ID:       ticket.ID,
		Customer: ticket.Customer,
	}

	ticketInterface, err := article.Interface()
	if err != nil {
		log.Error(err)
		return
	}

	data, err := t.client.TicketArticleCreate(ticketInterface)
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(bytes, &Article); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) Update(ticketID int, ticket model.BotTicket) (update *map[string]interface{}, err error) {
	ticketInterface, err := ticket.Interface()
	if err != nil {
		log.Error(err)
		return
	}

	if update, err = t.client.TicketUpdate(ticketID, ticketInterface); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) Delete(ticketID int) (err error) {
	if err = t.client.TicketDelete(ticketID); err != nil {
		log.Error(err)
		return err
	}

	return
}
