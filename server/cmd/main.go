package main

import (
	"log"

	"github.com/reyhanyogs/realtime-chat/db"
	"github.com/reyhanyogs/realtime-chat/internal/user"
	"github.com/reyhanyogs/realtime-chat/internal/ws"
	"github.com/reyhanyogs/realtime-chat/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not init database connection: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
