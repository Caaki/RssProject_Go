package main

import (
	"database/sql"
	"fmt"
	"github.com/Caaki/rssproject/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	feed, err := urlToFeed("https://wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)

	//=======Loading env variables=======
	godotenv.Load()
	portString := os.Getenv("PORT")
	dbConnection := os.Getenv("DB_URL")
	if portString == "" || dbConnection == "" {
		log.Fatal("Falied to load .emv")
	}
	fmt.Println("Port: " + portString)

	// ======== Connectiong to the database =========
	conn, err := sql.Open("postgres", dbConnection)

	if err != nil {
		fmt.Println("Error connecting to the database")
		log.Fatal(err)
	}

	queries := database.New(conn)

	db := database.New(conn)
	api := apiConfig{
		DB: queries,
	}
	//==========SCRAPING DATA==============

	go startScraping(db, 4, time.Minute)

	//========== Api configuration ========
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*,", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	router.Mount("/v1", v1Router)
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)

	v1Router.Post("/users", api.handlerCreateUser)
	v1Router.Get("/users", api.middlewareAuth(api.handlerGetUserByApiKey))

	v1Router.Post("/feeds", api.middlewareAuth(api.handlerCreateFeed))
	v1Router.Get("/feeds", api.handlerGetFeeds)

	v1Router.Post("/feed_follows", api.middlewareAuth(api.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", api.middlewareAuth(api.handlerGetFeedFollowsByUserId))
	v1Router.Get("/feed_follows_api", api.middlewareAuth(api.handlerGetFeedFollowsByUserApiKey))
	v1Router.Delete("/feed_follows/{feedFollowID}", api.middlewareAuth(api.handleDeleteFeedFollow))

	v1Router.Get("/post", api.middlewareAuth(api.handlerGetUserPosts))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
