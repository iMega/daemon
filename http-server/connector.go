// Copyright © 2020 Dmitry Stoletov <info@imega.ru>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpserver

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/imega/daemon"
	"github.com/improbable-eng/go-httpwares"
	http_logrus "github.com/improbable-eng/go-httpwares/logging/logrus"
	http_ctxtags "github.com/improbable-eng/go-httpwares/tags"
	"github.com/sirupsen/logrus"
)

// Connector is a wrapped http server.
type Connector struct {
	srv     *http.Server
	conf    *Config
	log     logrus.FieldLogger
	prefix  string
	handler http.Handler
	daemon.WatcherConfigFunc
	daemon.ShutdownFunc
}

// Config of connector.
type Config struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

const defaultTimeout = 2 * time.Second

// New get a instance of http server.
func New(prefix string, l logrus.FieldLogger, handler http.Handler) *Connector {
	c := &Connector{
		log: l,
		conf: &Config{
			Addr:         "0.0.0.0:65534",
			ReadTimeout:  defaultTimeout,
			WriteTimeout: defaultTimeout,
		},
		handler: handler,
		prefix:  prefix,
	}

	c.WatcherConfigFunc = func() daemon.WatcherConfig {
		return daemon.WatcherConfig{
			Prefix:    prefix,
			MainKey:   "http-server",
			Keys:      []string{""},
			ApplyFunc: c.connect,
		}
	}

	c.ShutdownFunc = func() {
		if c.srv == nil {
			return
		}

		if err := c.srv.Shutdown(context.Background()); err != nil {
			c.log.Error(err)
		}
	}

	return c
}

func (c *Connector) newServer() *http.Server {
	var log *logrus.Entry
	if e, ok := c.log.(*logrus.Entry); ok {
		log = e
	}

	opts := []http_logrus.Option{
		http_logrus.WithDecider(
			func(w httpwares.WrappedResponseWriter, r *http.Request) bool {
				return w.StatusCode() != http.StatusOK
			},
		),
	}

	return &http.Server{
		Addr: c.conf.Addr,
		Handler: http_ctxtags.Middleware("http")(
			http_logrus.Middleware(log, opts...)(c.handler),
		),
		ReadTimeout:       c.conf.ReadTimeout,
		ReadHeaderTimeout: c.conf.ReadHeaderTimeout,
		WriteTimeout:      c.conf.WriteTimeout,
		IdleTimeout:       c.conf.IdleTimeout,
		MaxHeaderBytes:    c.conf.MaxHeaderBytes,
	}
}

func (c *Connector) connect(conf, last map[string]string) {
	reset := c.reset(last)
	config := c.config(conf)

	if !reset && !config {
		c.log.Debug("http connector has same configuration")

		return
	}

	if c.srv != nil {
		c.log.Debugf("http connector start shutdown, %s", c.conf.Addr)

		if err := c.srv.Shutdown(context.Background()); err != nil {
			c.log.Error(err)
		}

		c.log.Debugf("http connector end shutdown")
		c.srv = nil
	}

	c.srv = c.newServer()

	go func() {
		c.log.Debugf("http connector start on %s", c.conf.Addr)

		if err := c.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				c.log.Error(err)
			}
		}
	}()
}

func (c *Connector) config(conf map[string]string) bool {
	needUpdate := false

	for k, v := range conf {
		switch k {
		case c.prefix + "/http-server/host":
			needUpdate = needUpdate || c.conf.Addr != v
			c.conf.Addr = v

		case c.prefix + "/http-server/read-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.conf.ReadTimeout != d
				c.conf.ReadTimeout = d
			}

		case c.prefix + "/http-server/read-header-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.conf.ReadHeaderTimeout != d
				c.conf.ReadHeaderTimeout = d
			}

		case c.prefix + "/http-server/write-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.conf.WriteTimeout != d
				c.conf.WriteTimeout = d
			}

		case c.prefix + "/http-server/idle-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.conf.IdleTimeout != d
				c.conf.IdleTimeout = d
			}

		case c.prefix + "/http-server/max-header-bytes":
			if i, err := strconv.Atoi(v); err == nil {
				needUpdate = needUpdate || c.conf.MaxHeaderBytes != i
				c.conf.MaxHeaderBytes = i
			}
		}
	}

	return needUpdate
}

func (c *Connector) reset(last map[string]string) bool {
	needUpdate := false

	for k := range last {
		switch k {
		case c.prefix + "/http-server/host":
			needUpdate = true
			c.conf.Addr = "0.0.0.0:65534"

		case c.prefix + "/http-server/read-timeout":
			needUpdate = true
			c.conf.ReadTimeout = defaultTimeout

		case c.prefix + "/http-server/read-header-timeout":
			needUpdate = true
			c.conf.ReadHeaderTimeout = defaultTimeout

		case c.prefix + "/http-server/write-timeout":
			needUpdate = true
			c.conf.WriteTimeout = defaultTimeout

		case c.prefix + "/http-server/idle-timeout":
			needUpdate = true
			c.conf.IdleTimeout = defaultTimeout

		case c.prefix + "/http-server/max-header-bytes":
			needUpdate = true
			c.conf.MaxHeaderBytes = http.DefaultMaxHeaderBytes
		}
	}

	return needUpdate
}
