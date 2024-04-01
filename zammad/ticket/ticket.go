package ticket

import (
	"encoding/json"

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

func (t *Ticket) TicketsByCustomer(customer int) (TicketsByCustomer model.TicketsByCustomer, err error) {
	list, err := t.client.TicketListByCustomer(customer)
	if err != nil {
		log.Error(err)
		return
	}

	listBytes, err := json.Marshal(*list)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(listBytes, &TicketsByCustomer); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) PriorityList() (Priorities []model.Priority, err error) {
	list, err := t.client.TicketPriorityList()
	if err != nil {
		log.Error(err)
		return
	}

	listBytes, err := json.Marshal(list)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(listBytes, &Priorities); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) Create(ticket model.Ticket) (TicketCreate model.TicketCreate, err error) {
	ticketInterface, err := ticket.Interface()
	if err != nil {
		log.Error(err)
		return
	}

	list, err := t.client.TicketCreate(ticketInterface)
	if err != nil {
		log.Error(err)
		return
	}

	listBytes, err := json.Marshal(list)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(listBytes, &TicketCreate); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) SendToTicket(ticket model.Ticket) (Article model.TicketArticle, err error) {
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

	list, err := t.client.TicketArticleCreate(ticketInterface)
	if err != nil {
		log.Error(err)
		return
	}

	listBytes, err := json.Marshal(list)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(listBytes, &Article); err != nil {
		log.Error(err)
		return
	}

	return
}

func (t *Ticket) Update(ticketID int, ticket model.Ticket) (update *map[string]interface{}, err error) {
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
