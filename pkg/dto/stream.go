package dto

type CreateStream struct {
	StreamWebsocketAddr string `json:"stream_websocket_addr"`
	ChatWebsocketAddr   string `json:"chat_websocket_addr"`
	ViewerWebsocketAddr string `json:"viewer_websocket_addr"`
	Type                string `json:"type"`
}

type CreateNoStream struct {
	NoStream string `json:"no_stream"`
	Leave    string `json:"leave"`
}
