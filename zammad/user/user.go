package user

import (
	"encoding/json"
	"errors"
	mail2 "net/mail"
	"slices"
	"strconv"

	"github.com/AlessandroSechi/zammad-go"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/client"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	log "github.com/sirupsen/logrus"
)

type User struct {
	zammad *zammad.Client
	client *client.Client
}

func NewUserController(zammad *zammad.Client, client *client.Client) *User {
	return &User{zammad: zammad, client: client}
}

func (u *User) Me(user, pass string) (User model.User, err error) {
	newZammad, err := u.client.NewClient(user, pass)
	if err != nil {
		log.Error(err)
		return
	}

	data, err := newZammad.UserMe()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal(bytes, &User); err != nil {
		log.Error(err)
	}

	return
}

func (u *User) ListByDepartment(department string, group int) (Users []model.User, err error) {
	data, err := u.zammad.UserList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var allUsers []model.User
	if err = json.Unmarshal(bytes, &allUsers); err != nil {
		log.Error(err)
		return
	}

	for _, user := range allUsers {
		if _, ok := user.GroupIds[strconv.Itoa(group)]; ok && user.Active && user.Department == department && slices.Contains(user.RoleIds, 2) {
			Users = append(Users, user)
		}
	}

	return
}

func (u *User) ValidateEmail(emails string) (err error) {
	if _, err = mail2.ParseAddress(emails); err != nil {
		log.Error(err)
		return
	}

	data, err := u.zammad.UserList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var allUsers []model.User
	if err = json.Unmarshal(bytes, &allUsers); err != nil {
		log.Error(err)
		return
	}

	for _, user := range allUsers {
		if user.Active && user.Email == emails {
			return
		}
	}

	return errors.New("not found")
}

func (u *User) Departments(group int) (departments []string, err error) {
	data, err := u.zammad.UserList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var allUsers []model.User
	if err = json.Unmarshal(bytes, &allUsers); err != nil {
		log.Error(err)
		return
	}

	for _, user := range allUsers {
		if _, ok := user.GroupIds[strconv.Itoa(group)]; ok && user.Active &&
			slices.Contains(user.RoleIds, 2) && user.Department != "" && !slices.Contains(departments, user.Department) {
			departments = append(departments, user.Department)
		}
	}

	return
}
