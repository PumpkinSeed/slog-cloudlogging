# slog: Cloud Logging handler

![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-%23007d9c)

A Google Cloud Logging Handler for [slog](https://pkg.go.dev/log/slog) Go library.

## Install

```sh
go get github.com/PumpkinSeed/slog-cloudlogging
```

## Usage

GoDoc: [https://pkg.go.dev/github.com/PumpkinSeed/slog-cloudlogging](https://pkg.go.dev/github.com/PumpkinSeed/slog-cloudlogging)

### Example

```go
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
```