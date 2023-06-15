package handlers

import (
	"net/http"
	"sync"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/pkg/chat"
	"github.com/fajarardiyanto/go-media-server/pkg/protocol"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	sync.Mutex
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{}
}

func (s *ChatHandler) RoomChat(c *gin.Context) {
	c.HTML(http.StatusOK, "layouts/main", nil)
}

func (s *ChatHandler) RoomChatWebsocket(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		return
	}

	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	s.Lock()
	room := protocol.Rooms[uuid]
	s.Unlock()
	if room == nil {
		config.GetLogger().Debug("room id is null")
		c.JSON(http.StatusBadRequest, model.Response{
			Error:   true,
			Message: "room id is null",
		})
		return
	}
	if room.Hub == nil {
		config.GetLogger().Debug("hub is nil")
		c.JSON(http.StatusBadRequest, model.Response{
			Error:   true,
			Message: "hub is null",
		})
		return
	}

	conn := chat.Client{}
	conn.PeerChatConn(unsafeConn, room.Hub)
}

func (s *ChatHandler) StreamChatWebsocket(c *gin.Context) {
	sid := c.Param("suuid")
	if sid == "" {
		config.GetLogger().Debug("stream id is nil")
		c.JSON(http.StatusBadRequest, model.Response{
			Error:   true,
			Message: "stream id is null",
		})
		return
	}

	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	s.Lock()
	if stream, ok := protocol.Streams[sid]; ok {
		s.Unlock()
		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		conn := chat.Client{}
		conn.PeerChatConn(unsafeConn, stream.Hub)
		return
	}
	s.Unlock()
}
