package group

import (
	"encoding/json"

	"github.com/AlessandroSechi/zammad-go"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	log "github.com/sirupsen/logrus"
)

type Group struct {
	client *zammad.Client
}

func NewGroupController(client *zammad.Client) *Group {
	return &Group{client: client}
}

func (t *Group) List() (Groups []model.Group, err error) {
	list, err := t.client.GroupList()
	if err != nil {
		log.Error(err)
		return
	}

	listBytes, err := json.Marshal(*list)
	if err != nil {
		log.Error(err)
		return
	}

	var allGroups []model.Group
	if err = json.Unmarshal(listBytes, &allGroups); err != nil {
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
