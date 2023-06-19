package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model/dto/response"
	"github.com/fajarardiyanto/go-media-server/pkg/chat"
	"github.com/fajarardiyanto/go-media-server/pkg/dto"
	"github.com/fajarardiyanto/go-media-server/pkg/protocol"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crypto/sha256"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type roomHandler struct {
	sync.Mutex
}

func NewRoomHandler() *roomHandler {
	return &roomHandler{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *roomHandler) RoomCreate(c *gin.Context) {
	id := uuid.New().String()
	c.JSON(http.StatusOK, response.Response{
		Error: false,
		Data:  id,
	})
}

func (s *roomHandler) Room(c *gin.Context) {
	id := c.Param("uuid")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Response{
			Error:   true,
			Message: "uuid can't be null",
		})
		return
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		ws = "wss"
	}

	id, sId, _ := s.createOrGetRoom(id)

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	data := dto.CreateRoom{
		RoomWebsocketAddr:   fmt.Sprintf("%s://%s/room/%s/websocket", ws, c.Request.Host, id),
		RoomLink:            fmt.Sprintf("%s://%s/room/%s", scheme, c.Request.Host, id),
		ChatWebsocketAddr:   fmt.Sprintf("%s://%s/room/%s/chat/websocket", ws, c.Request.Host, id),
		ViewerWebsocketAddr: fmt.Sprintf("%s://%s/room/%s/viewer/websocket", ws, c.Request.Host, id),
		StreamLink:          fmt.Sprintf("%s://%s/stream/%s", scheme, c.Request.Host, sId),
		Type:                "room",
	}

	c.JSON(http.StatusOK, response.Response{
		Error: false,
		Data:  data,
	})
}

func (s *roomHandler) RoomWebsocket(c *gin.Context) {
	id := c.Param("uuid")
	if id == "" {
		return
	}

	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		config.GetLogger().Error("upgrade: %v", err)
		c.JSON(http.StatusInternalServerError, response.Response{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	_, _, room := s.createOrGetRoom(id)
	protocol.RoomConn(unsafeConn, room.Peers)
}

func (s *roomHandler) createOrGetRoom(rId string) (string, string, *protocol.Room) {
	s.Lock()
	defer s.Unlock()

	h := sha256.New()
	h.Write([]byte(rId))
	id := fmt.Sprintf("%x", h.Sum(nil))

	if room := protocol.Rooms[rId]; room != nil {
		if _, ok := protocol.Streams[id]; !ok {
			protocol.Streams[id] = room
		}
		return rId, id, room
	}

	hub := chat.NewHub()
	p := &protocol.Peers{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &protocol.Room{
		Peers: p,
		Hub:   hub,
	}

	protocol.Rooms[rId] = room
	protocol.Streams[id] = room

	go hub.Run()
	return rId, id, room
}

func (s *roomHandler) RoomViewerWebsocket(c *gin.Context) {
	id := c.Param("uuid")
	if id == "" {
		return
	}

	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		config.GetLogger().Error("upgrade: %v", err)
		return
	}

	s.Lock()
	if peer, ok := protocol.Rooms[id]; ok {
		s.Unlock()
		s.roomViewerConn(unsafeConn, peer.Peers)
		return
	}
	s.Unlock()
}

func (s *roomHandler) roomViewerConn(c *websocket.Conn, p *protocol.Peers) {
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
