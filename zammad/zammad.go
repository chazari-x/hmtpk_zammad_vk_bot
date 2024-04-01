package zammad

import (
	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/client"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/group"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/object"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/ticket"
	"github.com/chazari-x/hmtpk_zammad_vk_bot/zammad/user"
	"github.com/chazari-x/zammad-go"
)

type Zammad struct {
	User   *user.User
	Group  *group.Group
	Ticket *ticket.Ticket
	Type   *object.Object
}

func NewZammad(cfg config.Zammad) (*Zammad, error) {
	c, err := zammad.NewClient(&zammad.Client{
		Token: cfg.Token,
		Url:   cfg.Url,
	})
	if err != nil {
		return nil, err
	}

	return &Zammad{
		User:   user.NewUserController(c, client.NewClientController(c, cfg.Url)),
		Group:  group.NewGroupController(c),
		Ticket: ticket.NewTicketController(c),
		Type:   object.NewObjectController(c),
	}, nil
}
