package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT is not found")
	}
	fmt.Println("PORT: ", portString)

	router := chi.NewRouter()

	// cors configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	// only fire on GET requests
	v1Router.Get("/healthz", HandlerReadiness)
	v1Router.Get("/err", HandlerError)
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server is running on port: %v", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
