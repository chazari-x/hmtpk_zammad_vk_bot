package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Priority struct {
	Name string `json:"name"`
}

type Object struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Active     bool       `json:"active,omitempty"`
	DataOption DataOption `json:"data_option"`
}

type DataOption struct {
	Options map[string]string `json:"options"`
}

type Group struct {
	ID     int    `json:"id"`
	Active bool   `json:"active,omitempty"`
	Name   string `json:"name"`
}

type User struct {
	ID          int                 `json:"id"`
	Firstname   string              `json:"firstname"`
	Lastname    string              `json:"lastname"`
	Email       string              `json:"email"`
	Department  string              `json:"department"`
	Active      bool                `json:"active"`
	RoleIds     []int               `json:"role_ids"`
	GroupIds    map[string][]string `json:"group_ids"`
	DisplayName string              `json:"displayname"`
}

type Ticket struct {
	ArticleCount          int       `json:"article_count"`
	ArticleIds            []int     `json:"article_ids"`
	CreateArticleSenderID int       `json:"create_article_sender_id"`
	CreateArticleTypeID   int       `json:"create_article_type_id"`
	CreatedAt             time.Time `json:"created_at"`
	CreatedByID           int       `json:"created_by_id"`
	CustomerID            int       `json:"customer_id"`
	GroupID               int       `json:"group_id"`
	ID                    int       `json:"id"`
	Number                string    `json:"number"`
	OwnerID               int       `json:"owner_id"`
	PriorityID            int       `json:"priority_id"`
	StateID               int       `json:"state_id"`
	Title                 string    `json:"title"`
	Type                  string    `json:"type"`
	UpdatedAt             time.Time `json:"updated_at"`
	UpdatedByID           int       `json:"updated_by_id"`
}

type BotTicket struct {
	ID         int     `json:"id"`
	Number     string  `json:"number"`
	Title      string  `json:"title"`
	Group      Group   `json:"group"`
	Customer   string  `json:"customer"`
	Priority   string  `json:"priority"`
	Department string  `json:"department"`
	Owner      Owner   `json:"owner"`
	Type       Type    `json:"type"`
	Article    Article `json:"article"`
	Status     string  `json:"status"`
}

type Owner struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Type struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Article struct {
	Body string `json:"body"`
}

type TicketArticleCreate struct {
	Body        string `json:"body"`
	ContentType string `json:"content_type"`
	Customer    string `json:"customer"`
	ID          int    `json:"id"`
	Sender      string `json:"sender"`
	Type        string `json:"type"`
}

type TicketArticle struct {
	Body        string    `json:"body"`
	ContentType string    `json:"content_type"`
	Customer    string    `json:"customer"`
	From        string    `json:"from"`
	ID          int       `json:"id"`
	Internal    bool      `json:"internal"`
	Sender      string    `json:"sender"`
	SenderID    int       `json:"sender_id"`
	TicketID    int       `json:"ticket_id"`
	Attachments []any     `json:"attachments"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	CreatedByID int       `json:"created_by_id"`
	Type        string    `json:"type"`
	TypeID      int       `json:"type_id"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedByID int       `json:"updated_by_id"`
}

func (t TicketArticleCreate) Interface() (*map[string]interface{}, error) {
	ticket := make(map[string]interface{})

	ticket["ticket_id"] = t.ID
	ticket["customer_id"] = t.Customer
	ticket["body"] = t.Body
	ticket["sender"] = "Customer"
	ticket["content_type"] = "text/plain"
	ticket["type"] = "phone"

	return &ticket, nil
}

const p = "\n# %s%s\n"

func (t BotTicket) String() (result string) {
	if t.Number != "" {
		if t.Title != "" {
			result += fmt.Sprintf(p, fmt.Sprintf("#%s \"%s\"", t.Number, t.Title), "")
		} else {
			result += fmt.Sprintf(p, "#", t.Number)
		}
	} else if t.Title != "" {
		result += fmt.Sprintf(p, "Заголовок: ", t.Title)
	}
	if t.Article.Body != "" {
		result += fmt.Sprintf(p, "Описание: ", t.Article.Body)
	}
	if t.Group.Name != "" {
		result += fmt.Sprintf(p, "Группа: ", t.Group.Name)
	}
	if t.Type.Value != "" {
		result += fmt.Sprintf(p, "Тип: ", t.Type.Value)
	}
	if t.Priority != "" {
		result += fmt.Sprintf(p, "Приоритет: ", t.Priority)
	}
	if t.Department != "" {
		result += fmt.Sprintf(p, "Отдел: ", t.Department)
	}
	if t.Owner.Name != "" {
		result += fmt.Sprintf(p, "Ответственный: ", t.Owner.Name)
	}
	if t.Status != "" {
		result += fmt.Sprintf(p, "Статус: ", t.Status)
	}

	return
}

func (t BotTicket) Interface() (*map[string]interface{}, error) {
	if t.Title == "" {
		return nil, TitleIsNil
	}

	if t.Article.Body == "" {
		return nil, BodyIsNil
	}

	if t.Group.Name == "" {
		return nil, GroupIsNil
	}

	customer, err := strconv.Atoi(t.Customer)
	if err != nil || customer == 0 {
		return nil, CustomerIsNil
	}

	ticket := make(map[string]interface{})

	if t.Priority != "" {
		ticket["priority"] = t.Priority
	}

	if t.Owner.ID != "" {
		ticket["owner_id"] = t.Owner.ID
	}

	if t.Type.Key != "" {
		ticket["type"] = t.Type.Key
	}

	ticket["customer_id"] = t.Customer
	ticket["title"] = t.Title
	ticket["group"] = t.Group.Name

	article := make(map[string]interface{})

	article["body"] = t.Article.Body
	article["sender"] = "Customer"
	article["content_type"] = "text/plain"
	article["type"] = "phone"

	ticket["article"] = article

	return &ticket, nil
}

var (
	TitleIsNil    = errors.New("title is nil")
	BodyIsNil     = errors.New("body is nil")
	GroupIsNil    = errors.New("group is nil")
	CustomerIsNil = errors.New("customer is nil")
)

type TicketsByCustomer struct {
	TicketIdsOpen   []int `json:"ticket_ids_open"`
	TicketIdsClosed []int `json:"ticket_ids_closed"`
	Assets          struct {
		Ticket map[string]TicketByCustomer `json:"Ticket"`
	} `json:"assets"`
}

type TicketByCustomer struct {
	ID         int    `json:"id"`
	GroupID    int    `json:"group_id"`
	PriorityID int    `json:"priority_id"`
	StateID    int    `json:"state_id"`
	Number     string `json:"number"`
	Title      string `json:"title"`
	OwnerID    int    `json:"owner_id"`
	CustomerID int    `json:"customer_id"`
	Type       any    `json:"type"`
}
