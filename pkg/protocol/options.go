package protocol

import (
	"sync"

	"github.com/fajarardiyanto/go-media-server/pkg/chat"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var (
	Rooms   map[string]*Room
	Streams map[string]*Room
)

var (
	turnConfig = webrtc.Configuration{
		ICETransportPolicy: webrtc.ICETransportPolicyRelay,
		ICEServers: []webrtc.ICEServer{
			{

				URLs: []string{"stun:turn.localhost:3478"},
			},
			{

				URLs: []string{"turn:turn.localhost:3478"},

				Username: "fajar",

				Credential:     "fajar",
				CredentialType: webrtc.ICECredentialTypePassword,
			},
		},
	}
)

type WebsocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Room struct {
	Peers *Peers
	Hub   *chat.Hub
}

type Peers struct {
	sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	Websocket      *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn *websocket.Conn
	sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteJSON(v)
}
