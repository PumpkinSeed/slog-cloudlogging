package main

import (
	"errors"
	"log/slog"
	"time"

	slogcloudlogging "github.com/PumpkinSeed/slog-cloudlogging"
)

func main() {
	googleHandler := slogcloudlogging.NewHandler("test", "test-logs", nil)
	googleHandler.AutoFlush()
	slog.SetDefault(slog.New(googleHandler))

	slog.Error("test message", slog.Any("error", errors.New("this is an error")))
	time.Sleep(2 * time.Second) // Wait for the Flush
}
