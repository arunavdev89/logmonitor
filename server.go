package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second
)

func prometheusHandler() http.Handler {
	return prometheus.Handler()
}

type HttpServer struct {
	Server *http.Server
}

//NewHttpServer creates a Http Server that exposes a /metrics endpoint
//Prometheus scrapes this endpoint and gather metrics
func NewHttpServer(cfg *Config) *HttpServer {
	router := mux.NewRouter()
	router.Handle("/metrics", prometheusHandler())
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Port),
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}
	return &HttpServer{
		Server: s,
	}
}

func (s *HttpServer) Serve() {
	s.Server.ListenAndServe()
}
