package server

import (
	"time"

	"github.com/fajarardiyanto/go-media-server/config"
	"github.com/fajarardiyanto/go-media-server/internal/model"
	"github.com/fajarardiyanto/go-media-server/internal/service"
	"github.com/fajarardiyanto/go-media-server/pkg/protocol"
	"github.com/gin-contrib/cors"

	"github.com/fajarardiyanto/go-media-server/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Run() error {
	config.Config()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
		AllowCredentials: true,
	}))

	// services
	userSvc := service.NewUserService()

	// handlers
	roomHandler := handlers.NewRoomHandler()
	chatHandler := handlers.NewChatHandler()
	streamHandler := handlers.NewStreamHandler()
	userHandler := handlers.NewUserHandler(userSvc)

	r.GET("/", handlers.Ping)
	r.POST("/register", userHandler.RegisterHandler)
	r.POST("/login", userHandler.LoginHandler)

	// r.Use(middleware.AuthMiddleware())

	r.GET("/room/create", roomHandler.RoomCreate)
	r.GET("/room/:uuid", roomHandler.Room)
	r.GET("/room/:uuid/websocket", roomHandler.RoomWebsocket)

	//r.GET("/room/:uuid/chat", chatHandler.RoomChat)
	r.GET("/room/:uuid/chat/websocket", chatHandler.RoomChatWebsocket)
	r.GET("/room/:uuid/viewer/websocket", roomHandler.RoomViewerWebsocket)

	r.GET("/stream/:suuid", streamHandler.Stream)
	r.GET("/stream/:suuid/websocket", streamHandler.StreamWebsocket)
	r.GET("/stream/:suuid/chat/websocket", chatHandler.StreamChatWebsocket)
	r.GET("/stream/:suuid/viewer/websocket", streamHandler.StreamViewerWebsocket)
	//r.Static("/static", "./assets")

	protocol.Rooms = make(map[string]*protocol.Room)
	protocol.Streams = make(map[string]*protocol.Room)

	go dispatchKeyFrames()
	return r.Run(":" + model.GetConfig().Port)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range protocol.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}
