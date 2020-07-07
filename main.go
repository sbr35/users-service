package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sbr35/wallets-users/db"
	"github.com/sbr35/wallets-users/handlers"
)

func main() {
	logger := log.New(os.Stdout, "users-api ", log.LstdFlags)
	collection, err := db.UsersCollection()
	if err != nil {
		log.Fatal(err)
	}
	usersHandler := handlers.NewUserHandler(logger, collection)
	loginHandler := handlers.NewLogin(logger, collection)
	ServeMux := http.NewServeMux()
	ServeMux.Handle("/api/v1/users", usersHandler)
	ServeMux.Handle("/api/v1/users/login", loginHandler)

	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      logRequest(ServeMux, logger),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		logger.Println("Server Started at port 8080")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	logger.Println("Received Terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(tc)
}

func logRequest(handler http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
