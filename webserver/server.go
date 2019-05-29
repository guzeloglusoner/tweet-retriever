package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// NewServer is used for instantinating socket server// NewServer creates a server instance
func NewServer(router *mux.Router, addr string) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
	})

	return &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      c.Handler(router),
		ErrorLog:     logger,
	}
}
