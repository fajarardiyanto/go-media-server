package protocol

import (
	"encoding/json"
	"github.com/fajarardiyanto/go-media-server/config"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

func (p *Peers) AddTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	p.Lock()
	defer func() {
		p.Unlock()
		p.SignalPeerConnections()
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		config.GetLogger().Error(err)
		return nil
	}

	p.TrackLocals[t.ID()] = trackLocal
	return trackLocal
}

func (p *Peers) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	p.Lock()
	defer func() {
		p.Unlock()
		p.SignalPeerConnections()
	}()

	delete(p.TrackLocals, t.ID())
}

func (p *Peers) SignalPeerConnections() {
	p.Lock()
	defer func() {
		p.Unlock()
		p.DispatchKeyFrame()
	}()

	attempt := func() (tryAgain bool) {
		for i := range p.Connections {
			if p.Connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				p.Connections = append(p.Connections[:i], p.Connections[i+1:]...)
				config.GetLogger().Info("connection %v", p.Connections)
				return true
			}

			existingSenders := map[string]bool{}
			for _, sender := range p.Connections[i].PeerConnection.GetSenders() {
				if sender.Track() == nil {
					config.GetLogger().Debug("track sender is nil")
					continue
				}

				existingSenders[sender.Track().ID()] = true

				if _, ok := p.TrackLocals[sender.Track().ID()]; !ok {
					if err := p.Connections[i].PeerConnection.RemoveTrack(sender); err != nil {
						config.GetLogger().Error(err)
						return true
					}
				}
			}

			for _, receiver := range p.Connections[i].PeerConnection.GetReceivers() {
				if receiver.Track() == nil {
					config.GetLogger().Debug("track receiver is nil")
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			for trackID := range p.TrackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := p.Connections[i].PeerConnection.AddTrack(p.TrackLocals[trackID]); err != nil {
						config.GetLogger().Error(err)
						return true
					}
				}
			}

			offer, err := p.Connections[i].PeerConnection.CreateOffer(nil)
			if err != nil {
				config.GetLogger().Error(err)
				return true
			}

			if err = p.Connections[i].PeerConnection.SetLocalDescription(offer); err != nil {
				config.GetLogger().Error(err)
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				config.GetLogger().Error(err)
				return true
			}

			if err = p.Connections[i].Websocket.WriteJSON(&WebsocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				config.GetLogger().Error(err)
				return true
			}
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			go func() {
				time.Sleep(time.Second * 3)
				p.SignalPeerConnections()
			}()
			return
		}

		if !attempt() {
			break
		}
	}
}

func (p *Peers) DispatchKeyFrame() {
	p.Lock()
	defer p.Unlock()

	for i := range p.Connections {
		for _, receiver := range p.Connections[i].PeerConnection.GetReceivers() {
			if receiver.Track() == nil {
				config.GetLogger().Debug("track receiver is nil")
				continue
			}

			_ = p.Connections[i].PeerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
