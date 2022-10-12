package mongo

import (
	"github.com/globalsign/mgo"
)

type Client struct {
	session *mgo.Session
}

func NewClient(connectionString string) (*Client, error) {
	session, err := mgo.Dial(connectionString)
	if err != nil {
		return nil, err
	}
	return &Client{session: session}, nil
}

func (c *Client) NewSession() *mgo.Session {
	newSession := c.session.Copy()
	newSession.SetMode(mgo.Strong, true)
	return newSession
}
