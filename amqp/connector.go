// Copyright © 2022 Dmitry Stoletov <info@imega.ru>
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

package amqp

import (
	"encoding/json"
	"errors"

	"github.com/imega/daemon"
	"github.com/imega/daemon/logging"
	"github.com/imega/mt"
)

var ErrConnectorNotReady = errors.New("connector isn't ready")

type Connector struct {
	log      logging.Logger
	conf     Config
	MT       mt.MassTransport
	prefix   string
	handlers map[string]mt.HandlerFunc

	daemon.WatcherConfigFunc
	daemon.HealthCheckFunc
	daemon.ShutdownFunc
}

func (conn *Connector) Call(
	name mt.ServiceName,
	req mt.Request,
	rep mt.ReplyFunc,
) error {
	if conn.MT == nil {
		return ErrConnectorNotReady
	}

	return conn.MT.Call(name, req, rep)
}

func (conn *Connector) Cast(name mt.ServiceName, req mt.Request) error {
	if conn.MT == nil {
		return ErrConnectorNotReady
	}

	return conn.MT.Cast(name, req)
}

func (conn *Connector) ConnectAndServe() error { return nil }

func (conn *Connector) Shutdown() error { return nil }

func (conn *Connector) HealthCheck() bool { return false }

func (conn *Connector) AddHandler(name mt.ServiceName, handler mt.HandlerFunc) {
	if conn.MT != nil {
		return
	}

	conn.handlers[string(name)] = handler
}

func (conn *Connector) connect(conf, last map[string]string) {
	reset := conn.reset(last)
	config := conn.config(conf)

	if !reset && !config {
		conn.log.Debugf("mt connector has same configuration")

		return
	}

	if conn.MT != nil {
		conn.log.Debugf("mt connector starts to shutdown, %s", conn.conf.DSN)

		if err := conn.MT.Shutdown(); err != nil {
			conn.log.Errorf("failed to shutdown server, %s", err)
		}

		conn.log.Debugf("mt connector is stopped")
		conn.MT = nil
	}

	conn.MT = mt.NewMT(
		mt.WithAMQP(conn.conf.DSN),
		mt.WithLogger(conn.log),
		mt.WithConfig(conn.conf.Config),
	)

	for name, handler := range conn.handlers {
		conn.MT.AddHandler(mt.ServiceName(name), handler)
	}

	if err := conn.MT.ConnectAndServe(); err != nil {
		conn.log.Errorf("failed to start MassTransport, %s", err)
	}
}

func (conn *Connector) config(conf map[string]string) bool {
	needUpdate := false

	for k, value := range conf {
		switch k {
		case conn.prefix + "/mt/dsn":
			needUpdate = needUpdate || conn.conf.DSN != value
			conn.conf.DSN = value

		case conn.prefix + "/mt/config":
			if newConf, err := mt.ParseConfig([]byte(value)); err == nil {
				confRaw, err := json.Marshal(conn.conf.Config)
				if err != nil {
					conn.log.Errorf("failed to marshal config, %s", err)

					return false
				}

				needUpdate = needUpdate || string(confRaw) != value
				conn.conf.Config = newConf
			}
		}
	}

	return needUpdate
}

func (conn *Connector) reset(last map[string]string) bool {
	needUpdate := false

	for k := range last {
		switch k {
		case conn.prefix + "/mt/dsn":
			needUpdate = true
			conn.conf.DSN = ""
		case conn.prefix + "/mt/config":
			needUpdate = true
			conn.conf.Config = mt.Config{}
		}
	}

	return needUpdate
}

// Config of connector.
type Config struct {
	DSN    string
	Config mt.Config
}

type Option func(*Connector)

func New(prefix string, opts ...Option) *Connector {
	conn := &Connector{
		prefix:   prefix,
		conf:     Config{},
		handlers: make(map[string]mt.HandlerFunc),
	}

	for _, opt := range opts {
		opt(conn)
	}

	conn.WatcherConfigFunc = func() daemon.WatcherConfig {
		return daemon.WatcherConfig{
			Prefix:    prefix,
			MainKey:   "mt",
			Keys:      keys(),
			ApplyFunc: conn.connect,
		}
	}

	conn.HealthCheckFunc = func() bool {
		if conn.MT == nil {
			return false
		}

		return conn.MT.HealthCheck()
	}

	conn.ShutdownFunc = func() {
		if conn.MT == nil {
			return
		}

		if err := conn.MT.Shutdown(); err != nil {
			conn.log.Errorf("failed to shutdown server, %s", err)
		}
	}

	return conn
}

func keys() []string {
	return []string{
		"dsn",
		"config",
	}
}

func WithLogger(l logging.Logger) Option {
	return func(c *Connector) {
		c.log = l
	}
}

func WithHandler(serviceName string, handler mt.HandlerFunc) Option {
	return func(c *Connector) {
		c.handlers[serviceName] = handler
	}
}
