package dto

type CreateRoom struct {
	RoomWebsocketAddr   string `json:"room_websocket_addr,omitempty"`
	RoomLink            string `json:"room_link,omitempty"`
	ChatWebsocketAddr   string `json:"chat_websocket_addr,omitempty"`
	ViewerWebsocketAddr string `json:"viewer_websocket_addr,omitempty"`
	StreamLink          string `json:"stream_link,omitempty"`
	Type                string `json:"type,omitempty"`
}
