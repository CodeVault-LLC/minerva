package api

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/codevault-llc/humblebrag-api/internal/api/routes"
	"github.com/go-redis/redis/v8"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func Start(db *gorm.DB, redis *redis.Client, cache *scs.SessionManager) {
	api := routes.SetupRouter(db)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(api)

	fmt.Println("Server started on port 3000")
	err := http.ListenAndServe(":3000", cache.LoadAndSave(handler))
	if err != nil {
		panic(err)
	}
}
