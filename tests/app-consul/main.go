package main

import (
	"context"
	"net/http"
	"time"

	"github.com/imega/daemon"
	"github.com/imega/daemon/configuring/consul"
	grpcserver "github.com/imega/daemon/grpc-server"
	health "github.com/imega/daemon/health/grpc"
	httpserver "github.com/imega/daemon/http-server"
	"github.com/imega/daemon/logging/wraplogrus"
	"github.com/imega/daemon/mysql"
	redis "github.com/imega/daemon/redis/sentinel"
	"github.com/improbable-eng/go-httpwares/logging/logrus/ctxlogrus"
	"google.golang.org/grpc"
)

const shutdownTimeout = 15 * time.Second

func main() {
	log := wraplogrus.New(wraplogrus.Config{
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
		}),
	)

	g1 := grpcserver.New(
		"testclient",
		grpcserver.WithLogger(log),
		grpcserver.WithServices(func(s *grpc.Server) {
			health.New(s)
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	h := httpserver.New(
		"myhttp",
		httpserver.WithLogger(log),
		httpserver.WithHandler(mux),
	)

	cr := consul.Watch(
		log,
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

	muxTest := http.NewServeMux()
	muxTest.Handle(
		"/test_mysql_reconnect_between_instances",
		testMysqlReconnectBetweenInstances(m),
	)
	muxTest.Handle(
		"/test_redis_reconnect_between_instances",
		testRedisReconnectBetweenInstances(r),
	)

	srvTest := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: muxTest,
	}

	go func() {
		if err := srvTest.ListenAndServe(); err != nil {
			log.Fatalf("failed to serve, %s", err)
		}
	}()

	d.RegisterShutdownFunc(
		h.ShutdownFunc,
		g.ShutdownFunc,
		g1.ShutdownFunc,
		m.ShutdownFunc,
		r.ShutdownFunc,
	)

	d.RegisterShutdownFunc(func() {
		if err := srvTest.Shutdown(context.Background()); err != nil {
			log.Error(err)
		}
	})

	log.Info("daemon is started")

	if err := d.Run(shutdownTimeout); err != nil {
		log.Errorf("failed to loop until shutdown: %s", err)
	}

	log.Info("daemon is stopped")
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func testMysqlReconnectBetweenInstances(c *mysql.Connector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			log   = ctxlogrus.Extract(r.Context())
			title []byte
		)

		if err := c.DB.QueryRow("select title from test.test").Scan(&title); err != nil {
			log.Error(err)
		}

		if _, err := w.Write(title); err != nil {
			log.Error(err)
		}
	})
}

func testRedisReconnectBetweenInstances(c *redis.Connector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := ctxlogrus.Extract(r.Context())

		res, err := c.DB.Info("server").Bytes()
		if err != nil {
			log.Error(err)
		}

		if _, err := w.Write(res); err != nil {
			log.Error(err)
		}
	})
}
