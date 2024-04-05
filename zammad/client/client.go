package client

import (
	"github.com/chazari-x/zammad-go"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	zammad *zammad.Client
	url    string
}

func NewClientController(zammad *zammad.Client, url string) *Client {
	return &Client{zammad: zammad, url: url}
}

func (c *Client) NewClient(token string) (client *zammad.Client, err error) {
	if client, err = zammad.NewClient(&zammad.Client{OAuth: token, Url: c.url}); err != nil {
		log.Error(err)
		return
	}

	return
}
