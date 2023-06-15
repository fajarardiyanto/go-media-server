package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/pkg/dto"
	"github.com/fajarardiyanto/go-media-server/pkg/protocol"
	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

type StreamHandler struct {
	sync.Mutex
}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{}
}

func (s *StreamHandler) Stream(c *gin.Context) {
	suuid := c.Param("suuid")
	if suuid == "" {
		c.Status(400)
		return
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}

	s.Lock()
	if _, ok := protocol.Streams[suuid]; ok {
		s.Unlock()
		data := dto.CreateStream{
			StreamWebsocketAddr: fmt.Sprintf("%s://%s/stream/%s/websocket", ws, c.Request.Host, suuid),
			ChatWebsocketAddr:   fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Request.Host, suuid),
			ViewerWebsocketAddr: fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Request.Host, suuid),
			Type:                "stream",
		}

		c.JSON(http.StatusOK, model.Response{
			Error: false,
			Data:  data,
		})
		return
	}
	s.Unlock()

	c.JSON(http.StatusOK, model.Response{
		Error: false,
		Data: dto.CreateNoStream{
			NoStream: "true",
			Leave:    "true",
		},
	})
}

func (s *StreamHandler) StreamWebsocket(c *gin.Context) {
	suuid := c.Param("suuid")
	if suuid == "" {
		return
	}

	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	s.Lock()
	if stream, ok := protocol.Streams[suuid]; ok {
		s.Unlock()

		protocol.StreamConn(unsafeConn, stream.Peers)
		return
	}
	s.Unlock()
}

func (s *StreamHandler) StreamViewerWebsocket(c *gin.Context) {
	suuid := c.Param("suuid")
	if suuid == "" {
		return
	}

	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	s.Lock()
	if stream, ok := protocol.Streams[suuid]; ok {
		s.Unlock()
		s.viewerConn(unsafeConn, stream.Peers)
		return
	}
	s.Unlock()
}

func (s *StreamHandler) viewerConn(c *websocket.Conn, p *protocol.Peers) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ticker.C:
			w, err := c.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}
