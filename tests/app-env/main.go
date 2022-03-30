package main

import (
	"net/http"
	"time"

	"github.com/imega/daemon"
	"github.com/imega/daemon/configuring/env"
	grpcserver "github.com/imega/daemon/grpc-server"
	health "github.com/imega/daemon/health/grpc"
	httpserver "github.com/imega/daemon/http-server"
	"github.com/imega/daemon/logging"
	"github.com/imega/daemon/mysql"
	redis "github.com/imega/daemon/redis/sentinel"
	"google.golang.org/grpc"
)

const shutdownTimeout = 15 * time.Second

func main() {
	log := logging.New(logging.Config{
		Channel: "ch",
		Level:   "debug",
	})

	m := mysql.New("instance", "client", log)
	r := redis.New("instance", "rclient", log)

	g := grpcserver.New(
		"client",
		grpcserver.WithLogger(log),
		grpcserver.WithServices(func(s *grpc.Server) {
			health.New(s, m.HealthCheckFunc, r.HealthCheckFunc)
		}))

	g1 := grpcserver.New(
		"testclient",
		grpcserver.WithLogger(log),
		grpcserver.WithServices(func(s *grpc.Server) { health.New(s) }),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	h := httpserver.New(
		"myhttp",
		httpserver.WithLogger(log),
		httpserver.WithHandler(mux),
	)

	cr := env.Once(
		g.WatcherConfigFunc,
		g1.WatcherConfigFunc,
		m.WatcherConfigFuncs[0],
		m.WatcherConfigFuncs[1],
		r.WatcherConfigFuncs[0],
		r.WatcherConfigFuncs[1],
		h.WatcherConfigFunc,
	)

	d, err := daemon.New(log, cr)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("daemon is started")

	if err := d.Run(shutdownTimeout); err != nil {
		log.Errorf("failed to loop until shutdown: %s", err)
	}

	log.Info("daemon is stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
