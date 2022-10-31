package masstransport

import (
	"encoding/json"

	"github.com/imega/daemon"
	"github.com/imega/daemon/logging"
	"github.com/imega/mt"
)

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

func (c *Connector) AddHandler(serviceName string, handler mt.HandlerFunc) {
	if c.MT != nil {
		return
	}

	c.handlers[serviceName] = handler
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

func (c *Connector) config(conf map[string]string) bool {
	needUpdate := false

	for k, value := range conf {
		switch k {
		case c.prefix + "/mt/dsn":
			needUpdate = needUpdate || c.conf.DSN != value
			c.conf.DSN = value

		case c.prefix + "/mt/config":
			if newConf, err := mt.ParseConfig([]byte(value)); err == nil {
				confRaw, err := json.Marshal(c.conf.Config)
				if err != nil {
					c.log.Errorf("failed to marshal config, %s", err)

					return false
				}

				needUpdate = needUpdate || string(confRaw) != value
				c.conf.Config = newConf
			}
		}
	}

	return needUpdate
}

func (c *Connector) reset(last map[string]string) bool {
	needUpdate := false

	for k := range last {
		switch k {
		case c.prefix + "/mt/dsn":
			needUpdate = true
			c.conf.DSN = ""
		case c.prefix + "/mt/config":
			needUpdate = true
			c.conf.Config = mt.Config{}
		}
	}

	return needUpdate
}
