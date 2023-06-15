package chat

import (
	"bytes"
	"time"

	"github.com/fajarardiyanto/go-media-server/config"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		config.GetLogger().Error(err)
		return
	}

	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				config.GetLogger().Error(err)
				return
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				config.GetLogger().Error(err)
				return
			}

			if !ok {
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				config.GetLogger().Error(err)
				return
			}

			if _, err = w.Write(message); err != nil {
				config.GetLogger().Error(err)
				return
			}

			n := len(c.Send)
			for i := 0; i < n; i++ {
				if _, err = w.Write(newline); err != nil {
					config.GetLogger().Error(err)
					return
				}

				if _, err = w.Write(<-c.Send); err != nil {
					config.GetLogger().Error(err)
					return
				}
			}

			if err = w.Close(); err != nil {
				config.GetLogger().Error(err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				config.GetLogger().Error(err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				config.GetLogger().Error(err)
				return
			}
		}
	}
}

func (c *Client) PeerChatConn(conn *websocket.Conn, hub *Hub) {
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.register <- client

	go client.WritePump()
	client.ReadPump()
}
