package group

import (
	"encoding/json"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	"github.com/chazari-x/zammad-go"
	log "github.com/sirupsen/logrus"
)

type Group struct {
	client *zammad.Client
}

func NewGroupController(client *zammad.Client) *Group {
	return &Group{client: client}
}

func (t *Group) List() (Groups []model.Group, err error) {
	data, err := t.client.GroupList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var allGroups []model.Group
	if err = json.Unmarshal(bytes, &allGroups); err != nil {
		log.Error(err)
		return
	}

	for _, group := range allGroups {
		if group.Active {
			Groups = append(Groups, group)
		}
	}

	return
}
