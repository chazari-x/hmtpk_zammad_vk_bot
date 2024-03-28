package client

import (
	"github.com/AlessandroSechi/zammad-go"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	zammad *zammad.Client
	url    string
}

func NewClientController(zammad *zammad.Client, url string) *Client {
	return &Client{zammad: zammad, url: url}
}

func (c *Client) NewClient(user, pass string) (client *zammad.Client, err error) {
	if client, err = zammad.NewClient(&zammad.Client{
		Username: user,
		Password: pass,
		Url:      c.url,
	}); err != nil {
		log.Error(err)
		return
	}

	return
}
