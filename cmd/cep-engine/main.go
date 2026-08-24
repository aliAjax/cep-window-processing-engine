package main

import (
	"github.com/example/cep-window-engine/internal/checkpoint"
	"github.com/example/cep-window-engine/internal/domain"
	"github.com/example/cep-window-engine/internal/runtime"
	httptransport "github.com/example/cep-window-engine/internal/transport/http"
	"log"
	"net/http"
	"os"
)

func main() {
	store := domain.NewStore()
	if path := os.Getenv("CEP_RECOVERY_CHECKPOINT"); path != "" {
		action, err := recoverAtStartup(store, path)
		if err != nil {
			log.Printf("checkpoint recovery action=%s: %v", action, err)
		}
	}
	eng := runtime.New(store)
	dir := os.Getenv("CEP_CHECKPOINT_DIR")
	if dir == "" {
		dir = "./data/checkpoints"
	}
	srv := httptransport.New(store, eng, checkpoint.New(dir))
	addr := os.Getenv("CEP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("cep engine listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}
