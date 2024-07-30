package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/walteranderson/url-shortener/internal/database"
	"github.com/walteranderson/url-shortener/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}

	db, err := database.NewDatabaseConnection()
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	repo := database.NewRepository(db)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	router := router.NewRouter(repo)

	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, router))
}
