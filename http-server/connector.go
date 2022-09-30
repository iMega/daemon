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
	"github.com/imega/daemon/logging"
)

// Connector is a wrapped http server.
type Connector struct {
	srv     *http.Server
	conf    *Config
	log     logging.Logger
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

type Option func(*Connector)

func WithLogger(l logging.Logger) Option {
	return func(c *Connector) {
		c.log = l
	}
}

func WithHandler(o http.Handler) Option {
	return func(c *Connector) {
		c.handler = o
	}
}

// New get a instance of http server.
func New(prefix string, opts ...Option) *Connector {
	conn := &Connector{
		conf: &Config{
			Addr:         "0.0.0.0:65534",
			ReadTimeout:  defaultTimeout,
			WriteTimeout: defaultTimeout,
		},
		prefix: prefix,
		log:    logging.GetNoopLog(),
	}

	for _, opt := range opts {
		opt(conn)
	}

	conn.WatcherConfigFunc = func() daemon.WatcherConfig {
		return daemon.WatcherConfig{
			Prefix:    prefix,
			MainKey:   "http-server",
			Keys:      keys(),
			ApplyFunc: conn.connect,
		}
	}

	conn.ShutdownFunc = func() {
		if conn.srv == nil {
			return
		}

		if err := conn.srv.Shutdown(context.Background()); err != nil {
			conn.log.Errorf("%s", err)
		}
	}

	return conn
}

func (c *Connector) newServer() *http.Server {
	return &http.Server{
		Addr:              c.conf.Addr,
		Handler:           c.handler,
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
		c.log.Debugf("http connector has same configuration")

		return
	}

	if c.srv != nil {
		c.log.Debugf("http connector start shutdown, %s", c.conf.Addr)

		if err := c.srv.Shutdown(context.Background()); err != nil {
			c.log.Errorf("%s", err)
		}

		c.log.Debugf("http connector end shutdown")
		c.srv = nil
	}

	c.srv = c.newServer()

	go func() {
		c.log.Debugf("http connector start on %s", c.conf.Addr)

		if err := c.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				c.log.Errorf("%s", err)
			}
		}
	}()
}

func keys() []string {
	return []string{
		"host",
		"read-timeout",
		"read-header-timeout",
		"write-timeout",
		"idle-timeout",
		"max-header-bytes",
	}
}

func (c *Connector) config(conf map[string]string) bool {
	needUpdate := false

	for key, value := range conf {
		switch key {
		case c.prefix + "/http-server/host":
			needUpdate = needUpdate || c.conf.Addr != value
			c.conf.Addr = value

		case c.prefix + "/http-server/read-timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || c.conf.ReadTimeout != d
				c.conf.ReadTimeout = d
			}

		case c.prefix + "/http-server/read-header-timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || c.conf.ReadHeaderTimeout != d
				c.conf.ReadHeaderTimeout = d
			}

		case c.prefix + "/http-server/write-timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || c.conf.WriteTimeout != d
				c.conf.WriteTimeout = d
			}

		case c.prefix + "/http-server/idle-timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || c.conf.IdleTimeout != d
				c.conf.IdleTimeout = d
			}

		case c.prefix + "/http-server/max-header-bytes":
			if i, err := strconv.Atoi(value); err == nil {
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
