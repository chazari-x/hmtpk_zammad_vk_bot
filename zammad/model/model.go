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
	ID         int                 `json:"id"`
	Firstname  string              `json:"firstname"`
	Lastname   string              `json:"lastname"`
	Email      string              `json:"email"`
	Department string              `json:"department"`
	Active     bool                `json:"active"`
	RoleIds    []int               `json:"role_ids"`
	GroupIds   map[string][]string `json:"group_ids"`
}

type TicketCreate struct {
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
	//CloseAt                   any       `json:"close_at"`
	//CloseDiffInMin            any       `json:"close_diff_in_min"`
	//CloseEscalationAt         any       `json:"close_escalation_at"`
	//CloseInMin                any       `json:"close_in_min"`
	//EscalationAt              any       `json:"escalation_at"`
	//FirstResponseAt           any       `json:"first_response_at"`
	//FirstResponseDiffInMin    any       `json:"first_response_diff_in_min"`
	//FirstResponseEscalationAt any       `json:"first_response_escalation_at"`
	//FirstResponseInMin        any       `json:"first_response_in_min"`
	//LastCloseAt               any       `json:"last_close_at"`
	//LastContactAgentAt        any       `json:"last_contact_agent_at"`
	//LastContactAt             any       `json:"last_contact_at"`
	//LastContactCustomerAt     any       `json:"last_contact_customer_at"`
	//LastOwnerUpdateAt         any       `json:"last_owner_update_at"`
	//Note                      any       `json:"note"`
	//OrganizationID            any       `json:"organization_id"`
	//PendingTime               any       `json:"pending_time"`
	//Preferences               struct {
	//} `json:"preferences"`
	//TicketTimeAccountingIds []any     `json:"ticket_time_accounting_ids"`
	//TimeUnit                any       `json:"time_unit"`
	//UpdateDiffInMin         any       `json:"update_diff_in_min"`
	//UpdateEscalationAt      any       `json:"update_escalation_at"`
	//UpdateInMin             any       `json:"update_in_min"`
}

type Ticket struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Group      Group   `json:"group"`
	Customer   string  `json:"customer"`
	Priority   string  `json:"priority"`
	Department string  `json:"department"`
	Owner      Owner   `json:"owner"`
	Type       Type    `json:"type"`
	Article    Article `json:"article"`
}

type Owner struct {
	Name string `json:"name"`
	ID   string `json:"email"`
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
	//Cc           any       `json:"cc"`
	//InReplyTo    any       `json:"in_reply_to"`
	//MessageID    any       `json:"message_id"`
	//MessageIDMd5 any       `json:"message_id_md5"`
	//OriginByID   any       `json:"origin_by_id"`
	//Preferences  struct {
	//} `json:"preferences"`
	//References  any       `json:"references"`
	//ReplyTo     any       `json:"reply_to"`
	//Subject     any       `json:"subject"`
	//TimeUnit    any       `json:"time_unit"`
	//To          any       `json:"to"`
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

func (t Ticket) String() string {
	result := ""

	if t.Title != "" {
		result += fmt.Sprintf("\n▪ Заголовок:\n%s\n", t.Title)
	}
	if t.Article.Body != "" {
		result += fmt.Sprintf("\n▪ Описание:\n%s\n", t.Article.Body)
	}
	if t.Group.Name != "" {
		result += fmt.Sprintf("\n▪ Группа:\n%s\n", t.Group.Name)
	}
	if t.Type.Value != "" {
		result += fmt.Sprintf("\n▪ Тип:\n%s\n", t.Type.Value)
	}
	if t.Priority != "" {
		result += fmt.Sprintf("\n▪ Приоритет:\n%s\n", t.Priority)
	}
	if t.Department != "" {
		result += fmt.Sprintf("\n▪ Отдел:\n%s\n", t.Department)
	}
	if t.Owner.Name != "" {
		result += fmt.Sprintf("\n▪ Ответственное лицо:\n%s\n", t.Owner.Name)
	}

	return result
}

func (t Ticket) Interface() (*map[string]interface{}, error) {
	if t.Title == "" {
		return nil, TitleIsNil
	}

	if t.Article.Body == "" {
		return nil, BodyIsNil
	}

	if t.Group.Name == "" {
		return nil, GroupIsNil
	}

	//if t.Type.Key == "" {
	//	return nil, TypeIsNil
	//}

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
	//article["internal"] = true
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