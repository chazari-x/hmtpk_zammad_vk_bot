package keyboard

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	database "github.com/chazari-x/hmtpk_zammad_vk_bot/db"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/domain/vk-bot/model"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/storage"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad"
	model2 "github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	log "github.com/sirupsen/logrus"
)

type Keyboard struct {
	vk *api.VK
	z  *zammad.Zammad
	s  *storage.Storage
	db *database.DB
}

type Data struct {
	Page       int
	Department string
	Group      int
	Customer   int
	Link       string
}

func NewKeyboardGetter(vk *api.VK, z *zammad.Zammad, s *storage.Storage, db *database.DB) *Keyboard {
	return &Keyboard{vk: vk, z: z, s: s, db: db}
}

func (k *Keyboard) GetKeyboard(command model.Command, data Data) (marshal []byte, err error) {
	if data.Page < 1 {
		data.Page = 1
	}

	var kbrd model.Keyboard
	switch command.String() {
	case model.Home.String():
		if kbrd, err = k.home(data); err != nil {
			return
		}
	case model.CreateTicket.String():
		if kbrd, err = k.createTicket(data); err != nil {
			return
		}
	case model.ChangeTitle.String(), model.ChangeBody.String():
		if kbrd, err = k.changeTitleOrBody(data); err != nil {
			return
		}
	//case model.ChangeType.String():
	//	if kbrd, err = k.changeType(data); err != nil {
	//		return
	//	}
	case model.ChangeGroup.String():
		if kbrd, err = k.changeGroup(data); err != nil {
			return
		}
	//case model.ChangePriority.String():
	//	if kbrd, err = k.changePriority(data); err != nil {
	//		return
	//	}
	case model.ChangeDepartment.String():
		if kbrd, err = k.changeDepartment(data); err != nil {
			return
		}
	case model.ChangeOwner.String():
		if kbrd, err = k.changeOwner(data); err != nil {
			return
		}
	case model.MyTickets.String():
		if kbrd, err = k.myTickets(data); err != nil {
			return
		}
	case model.Authorization.String(), model.DeleteAuth.String():
		if kbrd, err = k.authorization(data); err != nil {
			return
		}
	case model.SendMessage.String():
		if kbrd, err = k.sendMessage(data); err != nil {
			return
		}
	default:
		kbrd = model.Keyboard{Buttons: [][]model.Button{}}
	}

	marshal, err = json.Marshal(kbrd)
	if err != nil {
		log.Error(err)
	}

	return
}

func (k *Keyboard) home(_ Data) (model.Keyboard, error) {
	return model.Keyboard{
		Buttons: [][]model.Button{{
			{
				Color: model.Negative,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.DeleteAuth.Button("")},
					Label:   model.DeleteAuth.String(),
				},
			},
			{
				Color: model.Primary,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.MyTickets.Button("")},
					Label:   model.MyTickets.String(),
				},
			},
		}},
	}, nil
}

func (k *Keyboard) authorization(data Data) (model.Keyboard, error) {
	return model.Keyboard{
		Buttons: [][]model.Button{
			{{
				Action: model.Action{
					Type:    "open_link",
					Payload: model.Payload{Button: model.Connect.Button("")},
					Label:   model.Connect.String(),
					Link:    data.Link,
				},
			}},
			{{
				Color: model.Positive,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.Authorization.Button("")},
					Label:   model.Authorization.String(),
				},
			}},
		},
	}, nil
}

func (k *Keyboard) sendMessage(_ Data) (model.Keyboard, error) {
	return model.Keyboard{
		Buttons: [][]model.Button{{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.CancelSend.Button("")},
				Label:   model.CancelSend.String(),
			},
		}}},
	}, nil
}

func (k *Keyboard) myTickets(data Data) (kbrd model.Keyboard, err error) {
	kbrd = model.Keyboard{Buttons: [][]model.Button{}}

	ticketsByCustomer, err := k.z.Ticket.TicketsByCustomer(data.Customer)
	if err != nil {
		return
	}

	var list []model2.TicketByCustomer
	for _, ticket := range ticketsByCustomer.Assets.Ticket {
		list = append(list, ticket)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID > list[j].ID
	})

	if float64(data.Page) > math.Ceil(float64(len(list))/4) {
		data.Page = 1
	}

	var first = (data.Page - 1) * 4
	var last = data.Page * 4

	if last > len(list) {
		last = len(list)
	}

	for _, element := range list[first:last] {
		key := model.Button{
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.SendMessage.Button(fmt.Sprintf("%d-%d", element.ID, element.CustomerID))},
				Label: func() string {
					if len(element.Title) > 37 {
						return fmt.Sprintf("%s...", element.Title[:37])
					}

					return element.Title
				}(),
			},
		}
		if slices.Contains(ticketsByCustomer.TicketIdsOpen, element.ID) {
			key.Color = model.Primary
		}
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{key})
	}

	if data.Page > 1 {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.MyTickets.Button("-" + strconv.Itoa(data.Page-1))},
				Label:   "Назад",
			},
		}})
		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
			kbrd.Buttons[len(kbrd.Buttons)-1] = append(kbrd.Buttons[len(kbrd.Buttons)-1], model.Button{
				Color: model.Positive,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.MyTickets.Button("-" + strconv.Itoa(data.Page+1))},
					Label:   "Дальше",
				},
			})
		}
	} else if float64(data.Page) < math.Ceil(float64(len(list))/4) {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.MyTickets.Button("-" + strconv.Itoa(data.Page+1))},
				Label:   "Дальше",
			},
		}})
	}

	kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
		Color: model.Negative,
		Action: model.Action{
			Type:    "callback",
			Payload: model.Payload{Button: model.Home.Button("")},
			Label:   model.Home.String(),
		},
	}})

	return
}

func (k *Keyboard) changeDepartment(data Data) (kbrd model.Keyboard, err error) {
	kbrd = model.Keyboard{Buttons: [][]model.Button{}}

	list, err := k.z.User.Departments(data.Group)
	if err != nil {
		return
	}

	if float64(data.Page) > math.Ceil(float64(len(list))/4) {
		data.Page = 1
	}

	var first = (data.Page - 1) * 4
	var last = data.Page * 4

	if last > len(list) {
		last = len(list)
	}

	for _, element := range list[first:last] {
		button, err := toButton(element, element)
		if err != nil {
			return kbrd, err
		}
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Primary,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: button},
				Label:   element,
			},
		}})
	}

	if data.Page > 1 {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.ChangeDepartment.Button("-" + strconv.Itoa(data.Page-1))},
				Label:   "Назад",
			},
		}})
		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
			kbrd.Buttons[len(kbrd.Buttons)-1] = append(kbrd.Buttons[len(kbrd.Buttons)-1], model.Button{
				Color: model.Positive,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.ChangeDepartment.Button("-" + strconv.Itoa(data.Page+1))},
					Label:   "Дальше",
				},
			})
		}
	} else if float64(data.Page) < math.Ceil(float64(len(list))/4) {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.ChangeDepartment.Button("-" + strconv.Itoa(data.Page+1))},
				Label:   "Дальше",
			},
		}})
	}

	kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
		Color: model.Negative,
		Action: model.Action{
			Type:    "callback",
			Payload: model.Payload{Button: model.Cancel.Button("")},
			Label:   model.Cancel.String(),
		},
	}})

	return
}

func (k *Keyboard) changeOwner(data Data) (kbrd model.Keyboard, err error) {
	kbrd = model.Keyboard{Buttons: [][]model.Button{}}

	list, err := k.z.User.ListByDepartment(data.Department, data.Group)
	if err != nil {
		return
	}

	if float64(data.Page) > math.Ceil(float64(len(list))/4) {
		data.Page = 1
	}

	var first = (data.Page - 1) * 4
	var last = data.Page * 4

	if last > len(list) {
		last = len(list)
	}

	for _, element := range list[first:last] {
		username := func() string {
			if element.DisplayName != "" && element.DisplayName != "-" {
				return element.DisplayName
			}
			return fmt.Sprintf("%s %s", element.Firstname, element.Lastname)
		}()
		button, err := toButton(strconv.Itoa(element.ID), username)
		if err != nil {
			return kbrd, err
		}
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Primary,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: button},
				Label:   username,
			},
		}})
	}

	if data.Page > 1 {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.ChangeOwner.Button("-" + strconv.Itoa(data.Page-1))},
				Label:   "Назад",
			},
		}})
		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
			kbrd.Buttons[len(kbrd.Buttons)-1] = append(kbrd.Buttons[len(kbrd.Buttons)-1], model.Button{
				Color: model.Positive,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.ChangeOwner.Button("-" + strconv.Itoa(data.Page+1))},
					Label:   "Дальше",
				},
			})
		}
	} else if float64(data.Page) < math.Ceil(float64(len(list))/4) {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.ChangeOwner.Button("-" + strconv.Itoa(data.Page+1))},
				Label:   "Дальше",
			},
		}})
	}

	kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
		Color: model.Negative,
		Action: model.Action{
			Type:    "callback",
			Payload: model.Payload{Button: model.Cancel.Button("")},
			Label:   model.Cancel.String(),
		},
	}})

	return
}

//func (k *Keyboard) changePriority(data Data) (kbrd model.Keyboard, err error) {
//	kbrd = model.Keyboard{Buttons: [][]model.Button{}}
//
//	list, err := k.z.BotTicket.PriorityList()
//	if err != nil {
//		return
//	}
//
//	if float64(data.Page) > math.Ceil(float64(len(list))/4) {
//		data.Page = 1
//	}
//
//	var first = (data.Page - 1) * 4
//	var last = data.Page * 4
//
//	if last > len(list) {
//		last = len(list)
//	}
//
//	for _, priority := range list[first:last] {
//		button, err := toButton(priority.Name, priority.Name)
//		if err != nil {
//			return kbrd, err
//		}
//		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//			Color: model.Primary,
//			Action: model.Action{
//				Type:    "callback",
//				Payload: model.Payload{Button: button},
//				Label:   priority.Name,
//			},
//		}})
//	}
//
//	if data.Page > 1 {
//		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//			Color: model.Positive,
//			Action: model.Action{
//				Type:    "callback",
//				Payload: model.Payload{Button: model.ChangePriority.Button("-" + strconv.Itoa(data.Page-1))},
//				Label:   "Назад",
//			},
//		}})
//		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
//			kbrd.Buttons[len(kbrd.Buttons)-1] = append(kbrd.Buttons[len(kbrd.Buttons)-1], model.Button{
//				Color: model.Positive,
//				Action: model.Action{
//					Type:    "callback",
//					Payload: model.Payload{Button: model.ChangePriority.Button("-" + strconv.Itoa(data.Page+1))},
//					Label:   "Дальше",
//				},
//			})
//		}
//	} else {
//		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
//			kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//				Color: model.Positive,
//				Action: model.Action{
//					Type:    "callback",
//					Payload: model.Payload{Button: model.ChangePriority.Button("-" + strconv.Itoa(data.Page+1))},
//					Label:   "Дальше",
//				},
//			}})
//		}
//	}
//
//	kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//		Color: model.Negative,
//		Action: model.Action{
//			Type:    "callback",
//			Payload: model.Payload{Button: model.Cancel.Button("")},
//			Label:   model.Cancel.String(),
//		}},
//	})
//
//	return
//}

func (k *Keyboard) changeGroup(data Data) (kbrd model.Keyboard, err error) {
	kbrd = model.Keyboard{Buttons: [][]model.Button{}}

	list, err := k.z.Group.List()
	if err != nil {
		return
	}

	if float64(data.Page) > math.Ceil(float64(len(list))/4) {
		data.Page = 1
	}

	var first = (data.Page - 1) * 4
	var last = data.Page * 4

	if last > len(list) {
		last = len(list)
	}

	for _, group := range list[first:last] {
		button, err := toButton(strconv.Itoa(group.ID), group.Name)
		if err != nil {
			return kbrd, err
		}
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Primary,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: button},
				Label:   group.Name,
			},
		}})
	}
	if data.Page > 1 {
		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
			Color: model.Positive,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.ChangeGroup.Button("-" + strconv.Itoa(data.Page-1))},
				Label:   "Назад",
			},
		}})
		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
			kbrd.Buttons[len(kbrd.Buttons)-1] = append(kbrd.Buttons[len(kbrd.Buttons)-1], model.Button{
				Color: model.Positive,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.ChangeGroup.Button("-" + strconv.Itoa(data.Page+1))},
					Label:   "Дальше",
				},
			})
		}
	} else {
		if float64(data.Page) < math.Ceil(float64(len(list))/4) {
			kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
				Color: model.Positive,
				Action: model.Action{
					Type:    "callback",
					Payload: model.Payload{Button: model.ChangeGroup.Button("-" + strconv.Itoa(data.Page+1))},
					Label:   "Дальше",
				},
			}})
		}
	}

	kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
		Color: model.Negative,
		Action: model.Action{
			Type:    "callback",
			Payload: model.Payload{Button: model.Cancel.Button("")},
			Label:   model.Cancel.String(),
		},
	}})

	return
}

type Type struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Types []Type

func (t Types) Len() int           { return len(t) }
func (t Types) Less(i, j int) bool { return t[i].Value < t[j].Value }
func (t Types) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

//func (k *Keyboard) changeType(data Data) (kbrd model.Keyboard, err error) {
//	kbrd = model.Keyboard{Buttons: [][]model.Button{}}
//
//	list, err := k.z.Type.List()
//	if err != nil {
//		return
//	}
//
//	var types Types
//	for _, object := range list {
//		if object.Name == "type" {
//			for key := range object.DataOption.Options {
//				types = append(types, Type{Key: key, Value: object.DataOption.Options[key]})
//			}
//			break
//		}
//	}
//
//	sort.Sort(types)
//
//	if float64(data.Page) > math.Ceil(float64(len(types))/4) {
//		data.Page = 1
//	}
//
//	var first = (data.Page - 1) * 4
//	var last = data.Page * 4
//
//	if last > len(types) {
//		last = len(types)
//	}
//
//	for _, t := range types[first:last] {
//		button, err := toButton(t.Key, t.Value)
//		if err != nil {
//			return kbrd, err
//		}
//		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//			Color: model.Primary,
//			Action: model.Action{
//				Type:    "callback",
//				Payload: model.Payload{Button: button},
//				Label:   t.Value,
//			},
//		}})
//	}
//
//	if data.Page > 1 {
//		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//			Color: model.Positive,
//			Action: model.Action{
//				Type:    "callback",
//				Payload: model.Payload{Button: model.ChangeType.Button("-" + strconv.Itoa(data.Page-1))},
//				Label:   "Назад",
//			},
//		}})
//		if float64(data.Page) < math.Ceil(float64(len(types))/4) {
//			kbrd.Buttons[len(kbrd.Buttons)-1] = append(kbrd.Buttons[len(kbrd.Buttons)-1], model.Button{
//				Color: model.Positive,
//				Action: model.Action{
//					Type:    "callback",
//					Payload: model.Payload{Button: model.ChangeType.Button("-" + strconv.Itoa(data.Page+1))},
//					Label:   "Дальше",
//				},
//			})
//		}
//	} else if float64(data.Page) < math.Ceil(float64(len(types))/4) {
//		kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//			Color: model.Positive,
//			Action: model.Action{
//				Type:    "callback",
//				Payload: model.Payload{Button: model.ChangeType.Button("-" + strconv.Itoa(data.Page+1))},
//				Label:   "Дальше",
//			},
//		}})
//	}
//
//	kbrd.Buttons = append(kbrd.Buttons, []model.Button{{
//		Color: model.Negative,
//		Action: model.Action{
//			Type:    "callback",
//			Payload: model.Payload{Button: model.Cancel.Button("")},
//			Label:   model.Cancel.String(),
//		},
//	}})
//
//	return
//}

func (k *Keyboard) changeTitleOrBody(_ Data) (model.Keyboard, error) {
	return model.Keyboard{
		Buttons: [][]model.Button{{{
			Color: model.Negative,
			Action: model.Action{
				Type:    "callback",
				Payload: model.Payload{Button: model.Cancel.Button("")},
				Label:   model.Cancel.String(),
			},
		}}},
	}, nil
}

func (k *Keyboard) createTicket(_ Data) (model.Keyboard, error) {
	return model.Keyboard{
		Buttons: [][]model.Button{
			{
				{
					Color: model.Primary,
					Action: model.Action{
						Type:    "callback",
						Payload: model.Payload{Button: model.ChangeTitle.Button("")},
						Label:   model.ChangeTitle.String(),
					},
				},
				{
					Color: model.Primary,
					Action: model.Action{
						Type:    "callback",
						Payload: model.Payload{Button: model.ChangeBody.Button("")},
						Label:   model.ChangeBody.String(),
					},
				},
			},
			{
				//{
				//	Color: model.Primary,
				//	Action: model.Action{
				//		Type:    "callback",
				//		Payload: model.Payload{Button: model.ChangeType.Button("")},
				//		Label:   model.ChangeType.String(),
				//	},
				//},
				{
					Color: model.Primary,
					Action: model.Action{
						Type:    "callback",
						Payload: model.Payload{Button: model.ChangeGroup.Button("")},
						Label:   model.ChangeGroup.String(),
					},
				},
				{
					Color: model.Primary,
					Action: model.Action{
						Type:    "callback",
						Payload: model.Payload{Button: model.ChangeOwner.Button("")},
						Label:   model.ChangeOwner.String(),
					},
				},
			},
			//{
			//	//{
			//	//	Color: model.Primary,
			//	//	Action: model.Action{
			//	//		Type:    "callback",
			//	//		Payload: model.Payload{Button: model.ChangePriority.Button("")},
			//	//		Label:   model.ChangePriority.String(),
			//	//	},
			//	//},
			//},
			{
				{
					Color: model.Positive,
					Action: model.Action{
						Type:    "callback",
						Payload: model.Payload{Button: model.Send.Button("")},
						Label:   model.Send.String(),
					},
				},
				{
					Color: model.Negative,
					Action: model.Action{
						Type:    "callback",
						Payload: model.Payload{Button: model.Delete.Button("")},
						Label:   model.Delete.String(),
					},
				},
			},
		},
	}, nil
}

func toButton(key, value string) (string, error) {
	marshal, err := json.Marshal(model.MorePayload{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Error(err)
	}
	return string(marshal), err
}
