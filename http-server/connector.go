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
	"net/http"
	"time"

	"github.com/imega/daemon"
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

	c.WatcherConfigFunc = func() ([]string, daemon.ApplyConfigFunc) {
		keys := []string{c.prefix + "/http-server"}

		return keys, c.connect
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

func (c *Connector) connect(conf, last map[string]string) {}
