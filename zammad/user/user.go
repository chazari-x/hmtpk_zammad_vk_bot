package user

import (
	"encoding/json"
	"slices"
	"strconv"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/client"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	"github.com/chazari-x/zammad-go"
	log "github.com/sirupsen/logrus"
)

type User struct {
	zammad *zammad.Client
	client *client.Client
}

func NewUserController(zammad *zammad.Client, client *client.Client) *User {
	return &User{zammad: zammad, client: client}
}

func (u *User) Me(token string) (User model.User, err error) {
	newZammad, err := u.client.NewClient(token)
	if err != nil {
		log.Error(err)
		return
	}

	data, err := newZammad.UserMe()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(data)
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

	bytes, err := json.Marshal(data)
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

func (u *User) Departments(group int) (departments []string, err error) {
	data, err := u.zammad.UserList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(data)
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
