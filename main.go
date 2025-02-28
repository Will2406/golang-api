package main

import (
	"context"
	"golang-api/handlers"
	"golang-api/middleware"
	"golang-api/server"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(
		context.Background(),
		&server.Config{
			JWTSecret:   JWT_SECRET,
			Port:        PORT,
			DatabaseUrl: DATABASE_URL,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, router *mux.Router) {
	router.Use(middleware.CheckAuthMiddleware(s))
	router.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	router.HandleFunc("/signup", handlers.SignUpHanlder(s)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	router.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
}
