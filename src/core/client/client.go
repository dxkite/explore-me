package client

import (
	"strconv"

	"dxkite.cn/explore-me/src/core/utils"
	"dxkite.cn/log"
	"golang.org/x/net/websocket"
)

const (
	TypeClientCount = "websocket:clientCount"
)

type ClientIdGetter func(*websocket.Conn) string
type ClientHandler func(*ClientPool, string, Message, *Client) error

type Client struct {
	Conn *websocket.Conn
}

type ClientPool struct {
	clients        map[string]*Client
	GetClientId    ClientIdGetter
	HandlerMessage ClientHandler
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func NewClientPool() *ClientPool {
	return &ClientPool{clients: map[string]*Client{}}
}

func (c *ClientPool) HandleClient(conn *websocket.Conn) {
	clientId := c.GetClientId(conn)
	client := &Client{Conn: conn}
	c.clients[clientId] = client
	c.broadcastClientCount()
	for {
		msg := Message{}
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			log.Error("Receive", clientId, err)
			delete(c.clients, clientId)
			c.broadcastClientCount()
			conn.Close()
			return
		} else {
			c.handleMessage(clientId, msg, client)
		}
	}
}

func (c *ClientPool) broadcastClientCount() {
	c.Broadcast(&Message{Type: TypeClientCount, Data: strconv.Itoa(c.Len())})
}

func (c *ClientPool) handleMessage(id string, msg Message, cli *Client) error {
	if c.HandlerMessage != nil {
		return c.HandlerMessage(c, id, msg, cli)
	}
	return nil
}

func (c *ClientPool) getClientId(conn *websocket.Conn) string {
	if c.GetClientId != nil {
		return c.GetClientId(conn)
	}
	return utils.GetRemoteIp(conn.Request())
}

func (c *ClientPool) Len() int {
	return len(c.clients)
}

func (c *ClientPool) Broadcast(msg *Message) int {
	succ := 0
	for id, conn := range c.clients {
		if err := websocket.JSON.Send(conn.Conn, msg); err != nil {
			log.Error("SendError", id, msg, err)
		} else {
			succ++
		}
	}
	return succ
}
