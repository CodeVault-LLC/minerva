package main

import (
	"log"
	"os"
	"time"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/updater"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
		log.Println(err)
	}

	parseUrl := os.Getenv("DATABASE_URL")
	stripeKey := os.Getenv("STRIPE_SECRET")
	stripe.Key = stripeKey

	constants.InitDB(parseUrl)
	constants.InitDragonflyDB()
	constants.InitAuth()
	constants.InitSessionManager()
	constants.InitConfig()

	go updater.StartAutoUpdate(20 * time.Minute)
	controller.Start()
}
