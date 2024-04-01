package object

import (
	"encoding/json"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/model"
	"github.com/chazari-x/zammad-go"
	log "github.com/sirupsen/logrus"
)

type Object struct {
	client *zammad.Client
}

func NewObjectController(client *zammad.Client) *Object {
	return &Object{client: client}
}

func (o *Object) List() (Objects []model.Object, err error) {
	list, err := o.client.ObjectList()
	if err != nil {
		log.Error(err)
		return
	}

	listBytes, err := json.Marshal(*list)
	if err != nil {
		log.Error(err)
		return
	}

	var allObjects []model.Object
	if err = json.Unmarshal(listBytes, &allObjects); err != nil {
		log.Error(err)
		return
	}

	for _, object := range allObjects {
		if object.Active {
			Objects = append(Objects, object)
		}
	}

	return
}
