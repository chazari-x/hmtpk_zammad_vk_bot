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
	data, err := o.client.ObjectList()
	if err != nil {
		log.Error(err)
		return
	}

	bytes, err := json.Marshal(*data)
	if err != nil {
		log.Error(err)
		return
	}

	var allObjects []model.Object
	if err = json.Unmarshal(bytes, &allObjects); err != nil {
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
