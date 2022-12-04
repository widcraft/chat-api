package chat

import "github.com/gorilla/websocket"

type user struct {
}

// chat client
type client struct {
	userIdx uint
	conn    *websocket.Conn
}

func (c *client) GetUserIdx() uint {
	return c.userIdx
}

func (c *client) SendMessage(message interface{}) error {
	return c.conn.WriteJSON(message)
}

// chat message
type message struct {
	userIdx uint
	message string
}