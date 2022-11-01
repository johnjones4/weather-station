package main

import (
	"context"
	"main/api"
	"main/core"
	"main/store"
	"main/transformers"
	"net/http"
	"os"
)

func main() {
	tx := []core.Transformer{
		transformers.AddDate,
		transformers.AdjustTemperature,
		transformers.CalculateWindSpeed,
	}

	db, err := store.NewPGStore(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	r := api.New(db, tx)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), r)
	if err != nil {
		panic(err)
	}
}
