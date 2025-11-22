package main

import (
	"context"
	"main/api"
	"main/core"
	"main/store"
	"main/transformers"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewDevelopmentConfig()
	l, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer l.Sync()
	log := l.Sugar()

	tx := []core.Transformer{
		transformers.AddDate,
	}

	var db core.Store
	pg := os.Getenv("POSTGRES_URL")
	if pg != "" {
		db, err = store.NewPGStore(context.Background(), os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
	} else {
		db = store.NewMemoryStore()
	}

	r := api.New(db, tx, log)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), r)
	if err != nil {
		panic(err)
	}
}
