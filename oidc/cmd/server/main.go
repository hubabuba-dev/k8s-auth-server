package main

import (
	"gokube/internal/auth"
	"gokube/internal/config"
	"gokube/internal/handlers"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[OIDC] ")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config init failed reason:", err)
	}
	log.Println("Config ready")

	keyManager, s, err := auth.CreateKeyManager(cfg.CronConfig)
	if err != nil {
		log.Println(err)
	}
	log.Println("test")

	defer s.Shutdown()

	tokenHandler := handlers.NewTokenHandler(keyManager)

	handler := mux.NewRouter()

	handler.HandleFunc("/token", tokenHandler.HandleToken).Methods("POST")
	handler.HandleFunc("/keys", tokenHandler.HandleJWKS)
	handler.HandleFunc("/.well-known/openid-configuration", handlers.HandleDiscovery).Methods("GET")

	err = http.ListenAndServe(":8443", handler)
	if err != nil {
		log.Fatal(err)
	}

}
