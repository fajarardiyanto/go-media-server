package protocol

import (
	"encoding/json"
	"os"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func RoomConn(c *websocket.Conn, p *Peers) {
	var cfg webrtc.Configuration
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		cfg = turnConfig
	}
	peerConnection, err := webrtc.NewPeerConnection(cfg)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}
	defer peerConnection.Close()

	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err = peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			config.GetLogger().Error(err)
			return
		}
	}

	newPeer := PeerConnectionState{
		PeerConnection: peerConnection,
		Websocket: &ThreadSafeWriter{
			Conn: c,
		}}

	p.Lock()
	p.Connections = append(p.Connections, newPeer)
	p.Unlock()

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			config.GetLogger().Error("candidate is nil")
			return
		}

		candidateString, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			config.GetLogger().Error(err)
			return
		}

		if err = newPeer.Websocket.WriteJSON(&WebsocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); err != nil {
			config.GetLogger().Error(err)
		}
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		switch state {
		case webrtc.PeerConnectionStateFailed:
			if err = peerConnection.Close(); err != nil {
				config.GetLogger().Error(err)
			}
		case webrtc.PeerConnectionStateClosed:
			p.SignalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(remote *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		trackLocal := p.AddTrack(remote)
		if trackLocal == nil {
			config.GetLogger().Error("track local is nil")
			return
		}
		defer p.RemoveTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := remote.Read(buf)
			if err != nil {
				config.GetLogger().Error(err)
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				config.GetLogger().Error(err)
				return
			}
		}
	})

	p.SignalPeerConnections()
	message := &WebsocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			config.GetLogger().Error(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			config.GetLogger().Error(err)
			return
		}

		switch message.Event {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err = json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				config.GetLogger().Error(err)
				return
			}

			if err = peerConnection.AddICECandidate(candidate); err != nil {
				config.GetLogger().Error(err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				config.GetLogger().Error(err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				config.GetLogger().Error(err)
				return
			}
		}
	}
}
