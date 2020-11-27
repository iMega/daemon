package cmd

func tplMain() string {
	return `package main

import (
    "errors"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/imega/daemon"
    "github.com/imega/daemon/logger"
)

const shutdownTimeout = 15 * time.Second

func main() {
    lconf := logger.Config{
        {{ if .Logger.Channel }}Channel: "{{ .Logger.Channel }}",{{ end }}
        {{ if .Logger.BuildID }}BuildID: "{{ .Logger.BuildID }}",{{ end }}
    }
    log := logger.New(lconf)
    d := daemon.New(log)
    if err := d.Run(shutdownTimeout); err != nil {
        log.Errorf("%w", err)
    }
}
`
}
