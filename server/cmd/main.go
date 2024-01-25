package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/reyhanyogs/realtime-chat/db"
	"github.com/reyhanyogs/realtime-chat/internal/user"
	"github.com/reyhanyogs/realtime-chat/internal/ws"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not init database connection: %s", err)
	}

	router := gin.Default()
	setupCORS(router)

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	user.NewHandler(router, userSvc)

	hub := ws.NewHub()
	ws.NewHandler(router, hub)
	go hub.Run()

	router.Run("0.0.0.0:8080")
}

func setupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))
}
