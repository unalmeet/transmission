package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	h  "ms/transmission/api"
	mr "ms/transmission/repository/cassandra"

	"ms/transmission/core"
)

func main() {
	repo := loadRepo()
	service := core.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/session/{idMeeting}", handler.GetSession)
	r.Post("/api/v1/session", handler.PostSession)
	r.Delete("/api/v1/session/{idMeeting}/{idSession}", handler.DeleteSession)
	r.Put("/api/v1/image", handler.PutImage)
	r.Put("/api/v1/sound", handler.PutAudio)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func loadRepo() core.ClientRepository {
	cassandraURL := os.Getenv("CASS_URL")
	cassandraDB := os.Getenv("CASS_DB")
	cassandraPASS := os.Getenv("CASS_PASS")
	cassandraUSR := os.Getenv("CASS_USR")
	repo, err := mr.NewRepository(cassandraURL, cassandraDB, cassandraUSR, cassandraPASS)
	if err != nil {
		log.Fatal(err)
	}
	return repo
}
