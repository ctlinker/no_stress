package main

import (
	"log"
	"net/http"
	"server/internal/api"
	"server/internal/db"
	"server/internal/env"
)

func main() {

	cfg := env.Load()
	database := db.Connect(cfg.DATABASE_URL)
	r := api.NewRouter(database)
	log.Println("listening on", cfg.PORT)
	http.ListenAndServe(":"+cfg.PORT, r)
}
