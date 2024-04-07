package model

import "fmt"

type Keyboard struct {
	// Скрывать ли клавиатуру после нажатия на клавишу, отправляющую сообщение. Например: text или location. Работает только для "inline": false.
	OneTime bool `json:"one_time,omitempty"`
	// True — клавиатура отображается внутри сообщения. False — показывает клавиатуру в диалоге, под полем ввода сообщения.
	Inline bool `json:"inline,omitempty"`
	// Массив строк с массивом клавиш.
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	// Объект, описывающий тип действия и его параметры.
	Action Action `json:"action"`
	// Цвет кнопки. Параметр используется только для кнопок типа text и callback.
	Color Color `json:"color,omitempty"`
}

type Color string

const (
	Primary   Color = "primary"
	Secondary Color = "secondary"
	Negative  Color = "negative"
	Positive  Color = "positive"
)

// Action : https://dev.vk.com/ru/api/bots/development/keyboard
type Action struct {
	Type    string  `json:"type"`
	AppID   string  `json:"app_id,omitempty"`
	OwnerID string  `json:"owner_id,omitempty"`
	Hash    string  `json:"hash,omitempty"`
	Label   string  `json:"label,omitempty"`
	Link    string  `json:"link,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Payload struct {
	Button string `json:"button,omitempty"`
}

type MorePayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ButtonPayload struct {
	Button MorePayload `json:"button,omitempty"`
}

type Message struct {
	MessagePayload string        `json:"-"`
	ButtonPayload  ButtonPayload `json:"-"`
	PeerID         int           `json:"-"`
	Text           string        `json:"-"`
	ID             int           `json:"-"`
}

const Status = "status"

type Command string

const (
	Authorization Command = "Начать"
	Connect       Command = "Связать"
	DeleteAuth    Command = "Выйти из системы"

	Home         Command = "На главную"
	MyTickets    Command = "Мои обращения"
	CreateTicket Command = "Создать обращение"

	ChangeTitle      Command = "Заголовок"
	ChangeBody       Command = "Описание"
	ChangeGroup      Command = "Группа"
	ChangeDepartment Command = "Отдел"
	ChangeSubject    Command = "Тема"
	ChangeOwner      Command = "Ответственный"
	ChangePriority   Command = "Приоритет"
	ChangeType       Command = "Тип"

	SendMessage Command = "Отправить ответ"

	Cancel     Command = "Отмена"
	CancelSend Command = "Отмена отправки"
	Delete     Command = "Удалить"
	Send       Command = "Отправить"
)

func (c Command) String() string {
	return string(c)
}

func (c Command) Button(n string) (str string) {
	return fmt.Sprintf("{\"key\": \"%s\", \"value\": \"%s\"}", c.Key()+n, c.Value()+n)
}

func (c Command) Key() string {
	switch c {
	case SendMessage:
		return "SendMessage"
	case CancelSend:
		return "CancelSend"
	case Connect:
		return "Connect"
	case MyTickets:
		return "MyTickets"
	case Home:
		return "Home"
	case CreateTicket:
		return "CreateTicket"
	case DeleteAuth:
		return "DeleteAuth"
	case Delete:
		return "Delete"
	case Cancel:
		return "Cancel"
	case ChangeOwner:
		return "ChangeOwner"
	case ChangeDepartment:
		return "ChangeDepartment"
	case Send:
		return "Send"
	//case ChangePriority:
	//	return "ChangePriority"
	case ChangeBody:
		return "ChangeBody"
	case ChangeTitle:
		return "ChangeTitle"
	case ChangeGroup:
		return "ChangeGroup"
	//case ChangeType:
	//	return "ChangeType"
	case ChangeSubject:
		return "ChangeSubject"
	case Authorization:
		return "Authorization"
	default:
		return c.String()
	}
}

func (c Command) Value() string {
	switch c {
	default:
		return c.String()
	}
}

func (c Command) Message() string {
	switch c {
	case Authorization:
		return `ℹ Для продолжения работы с ботом требуется связать VK с системой Zammad. Для этого нажмите на кнопку "Связать", затем нажмите кнопку "Начать".`
	case Home:
		return `ℹ При возникновении обращения просто напишите его мне (максимум 500 символов).`
	case DeleteAuth:
		return "ℹ Вы вышли из системы! \n" + Authorization.Message()
	case MyTickets:
		return "ℹ Ваши обращения:"
	case CreateTicket, Cancel:
		return "📄 Ваше обращение 📄\n"
	case ChangeTitle:
		return "➕ Введите заголовок (максимум 100 символов):\n"
	case ChangeBody:
		return "➕ Введите описание (максимум 500 символов):\n"
	case ChangeGroup:
		return "➕ Выберите группу:\n"
	//case ChangeType:
	//	return "➕ Выберите тип:\n"
	//case ChangePriority:
	//	return "➕ Выберите приоритет:\n"
	case ChangeDepartment:
		return "➕ Выберите отдел:\n"
	case ChangeSubject:
		return "➕ Введите тему:\n"
	case ChangeOwner:
		return "➕ Выберите ответственного:\n"
	case SendMessage:
		return "➕ Введите ваше сообщение:\n"
	case Send:
		return Send.String()
	default:
		return ""
	}
}

type Data struct {
	WhMsg   WebHookMessage
	Title   string
	Message string
	Kbrd    []byte
	Vk      int
}

type WebHookMessage struct {
	Ticket struct {
		CloseAt               any    `json:"close_at"`
		CloseDiffInMin        any    `json:"close_diff_in_min"`
		CloseEscalationAt     any    `json:"close_escalation_at"`
		CloseInMin            any    `json:"close_in_min"`
		CreateArticleSender   string `json:"create_article_sender"`
		CreateArticleSenderID int    `json:"create_article_sender_id"`
		CreateArticleType     string `json:"create_article_type"`
		CreateArticleTypeID   int    `json:"create_article_type_id"`
		CreatedAt             string `json:"created_at"`
		CustomerID            int    `json:"customer_id"`
		EscalationAt          any    `json:"escalation_at"`
		Group                 struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			UpdatedBy   string `json:"updated_by"`
			UpdatedByID int    `json:"updated_by_id"`
		} `json:"group"`
		ID                    int    `json:"id"`
		LastCloseAt           any    `json:"last_close_at"`
		LastContactAgentAt    any    `json:"last_contact_agent_at"`
		LastContactAt         string `json:"last_contact_at"`
		LastContactCustomerAt string `json:"last_contact_customer_at"`
		LastOwnerUpdateAt     any    `json:"last_owner_update_at"`
		Number                string `json:"number"`
		Owner                 struct {
			Displayname any    `json:"displayname"`
			Email       string `json:"email"`
			Firstname   string `json:"firstname"`
			ID          int    `json:"id"`
			Lastname    string `json:"lastname"`
			Mobile      string `json:"mobile"`
			Phone       string `json:"phone"`
		} `json:"owner"`
		Priority struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			UpdatedBy   string `json:"updated_by"`
			UpdatedByID int    `json:"updated_by_id"`
		} `json:"priority"`
		State     string `json:"state"`
		StateID   int    `json:"state_id"`
		TimeUnit  any    `json:"time_unit"`
		Title     string `json:"title"`
		Type      any    `json:"type"`
		UpdatedBy struct {
			Displayname string `json:"displayname"`
			Firstname   string `json:"firstname"`
			ID          int    `json:"id"`
			Lastname    string `json:"lastname"`
		} `json:"updated_by"`
	} `json:"ticket"`
	Article struct {
		Body      string `json:"body"`
		CreatedBy struct {
			Displayname string `json:"displayname"`
			Firstname   string `json:"firstname"`
			ID          int    `json:"id"`
			Lastname    string `json:"lastname"`
		} `json:"created_by"`
		Internal  bool `json:"internal"`
		UpdatedBy struct {
			Displayname string `json:"displayname"`
			Firstname   string `json:"firstname"`
			ID          int    `json:"id"`
			Lastname    string `json:"lastname"`
		} `json:"updated_by"`
	} `json:"article"`
}
