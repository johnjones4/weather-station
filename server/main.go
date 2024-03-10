package main

import (
	"context"
	"main/api"
	"main/core"
	"main/store"
	"main/transformers"
	"net/http"
	"os"
	"strconv"

	"github.com/johnjones4/errorbot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	chatId, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		panic(err)
	}
	bot := errorbot.New(
		"weather",
		os.Getenv("TELEGRAM_TOKEN"),
		chatId,
	)

	config := zap.NewDevelopmentConfig()
	l, err := config.Build(zap.Hooks(bot.ZapHook([]zapcore.Level{
		zapcore.FatalLevel,
		zapcore.PanicLevel,
		zapcore.DPanicLevel,
		zapcore.ErrorLevel,
		zapcore.WarnLevel,
	})))
	if err != nil {
		panic(err)
	}

	defer l.Sync()
	log := l.Sugar()

	tx := []core.Transformer{
		transformers.AddDate,
		transformers.AdjustTemperature,
	}

	db, err := store.NewPGStore(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	r := api.New(db, tx, log)
	err = http.ListenAndServe(os.Getenv("HTTP_HOST"), r)
	if err != nil {
		panic(err)
	}
}
