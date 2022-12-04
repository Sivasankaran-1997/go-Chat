package app

import (
	"chats/domain"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	r = gin.Default()
)

func StartApp() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")

	Routers()
	domain.ConnDB()
	r.Run(PORT)
}
