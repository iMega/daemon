package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/imega/daemon"
	"github.com/sirupsen/logrus"
)

// Connector is a wrapped redis sentinel client.
type Connector struct {
	log logrus.FieldLogger

	DB redis.UniversalClient

	opts    *redis.FailoverOptions
	pHost   string
	pClient string

	WatcherConfigFuncs []daemon.WatcherConfigFunc
	daemon.ShutdownFunc
	daemon.HealthCheckFunc
}

// New get a instance of redis sentinel client.
func New(pHost, pClient string, l logrus.FieldLogger) *Connector {
	c := &Connector{
		opts:    &redis.FailoverOptions{},
		pHost:   pHost + "/redis-sentinel/host",
		pClient: pClient + "/redis-sentinel",
		log:     l,
		DB:      &faker{},
	}

	if e, ok := l.(*logrus.Entry); ok {
		newLogger(e)
	}

	c.WatcherConfigFuncs = []daemon.WatcherConfigFunc{
		daemon.WatcherConfigFunc(func() daemon.WatcherConfig {
			return daemon.WatcherConfig{
				Prefix:    pHost,
				MainKey:   "redis-sentinel",
				Keys:      []string{"host"},
				ApplyFunc: c.connect,
			}
		}),
		daemon.WatcherConfigFunc(func() daemon.WatcherConfig {
			return daemon.WatcherConfig{
				Prefix:    pClient,
				MainKey:   "redis-sentinel",
				Keys:      clientConfig(),
				ApplyFunc: c.connect,
			}
		}),
	}

	c.ShutdownFunc = func() {
		if c.DB == nil {
			c.log.Error("failed to close connection to redis")

			return
		}

		if err := c.DB.Close(); err != nil {
			c.log.Errorf("failed to close connection to redis, %s", err)
		}
	}

	c.HealthCheckFunc = func() bool {
		if c.DB == nil {
			c.log.Error("failed to ping redis")

			return false
		}

		if _, err := c.DB.Ping().Result(); err != nil {
			c.log.Error(err)

			return false
		}

		c.log.Debug("redis ping is ok")

		return true
	}

	return c
}

func (c *Connector) connect(conf, last map[string]string) {

	fmt.Printf("\n\n=======\n%#v\n\n%#v\n\n", conf, last)

	reset := c.reset(last)
	config := c.config(conf)
	if !reset && !config {
		c.log.Debug("redis connector has same configuration")

		return
	}

	if _, ok := c.DB.(*faker); !ok {
		if err := c.DB.Close(); err != nil {
			c.log.Error(err)
		}
		c.log.Debug("redis connection closed")
		c.DB = &faker{}
	}

	c.DB = redis.NewFailoverClient(c.opts)

	c.log.Debug("redis connection open")
}

func (c *Connector) config(conf map[string]string) bool {
	needUpdate := false
	reconfigure := false

	for k := range conf {
		if strings.HasPrefix(k, c.pHost) {
			reconfigure = true
			c.opts.SentinelAddrs = nil
		}
	}

	if reconfigure {
		for k, v := range conf {
			if strings.HasPrefix(k, c.pHost) {
				if len(c.opts.SentinelAddrs) == 0 {
					needUpdate = true
					c.opts.SentinelAddrs = append(c.opts.SentinelAddrs, v)
				}

				for _, a := range c.opts.SentinelAddrs {
					if a == v {
						continue
					}
					needUpdate = true
					c.opts.SentinelAddrs = append(c.opts.SentinelAddrs, v)
				}
			}
		}
	}

	for k, v := range conf {
		switch k {
		case c.pClient + "/master-name":
			needUpdate = needUpdate || c.opts.MasterName != v
			c.opts.MasterName = v

		case c.pClient + "/password":
			needUpdate = needUpdate || c.opts.Password != v
			c.opts.Password = v

		case c.pClient + "/db":
			if i, err := strconv.Atoi(v); err == nil {
				needUpdate = needUpdate || c.opts.DB != i
				c.opts.DB = i
			}

		case c.pClient + "/max-retries":
			if i, err := strconv.Atoi(v); err == nil {
				needUpdate = needUpdate || c.opts.MaxRetries != i
				c.opts.MaxRetries = i
			}

		case c.pClient + "/min-retry-backoff":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.MinRetryBackoff != d
				c.opts.MinRetryBackoff = d
			}

		case c.pClient + "/max-retry-backoff":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.MaxRetryBackoff != d
				c.opts.MaxRetryBackoff = d
			}

		case c.pClient + "/dial-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.DialTimeout != d
				c.opts.DialTimeout = d
			}

		case c.pClient + "/read-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.ReadTimeout != d
				c.opts.ReadTimeout = d
			}

		case c.pClient + "/write-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.WriteTimeout != d
				c.opts.WriteTimeout = d
			}

		case c.pClient + "/pool-size":
			if i, err := strconv.Atoi(v); err == nil {
				needUpdate = needUpdate || c.opts.PoolSize != i
				c.opts.PoolSize = i
			}

		case c.pClient + "/min-idle-conns":
			if i, err := strconv.Atoi(v); err == nil {
				needUpdate = needUpdate || c.opts.MinIdleConns != i
				c.opts.MinIdleConns = i
			}

		case c.pClient + "/max-conn-age":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.MaxConnAge != d
				c.opts.MaxConnAge = d
			}

		case c.pClient + "/pool-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.PoolTimeout != d
				c.opts.PoolTimeout = d
			}

		case c.pClient + "/idle-timeout":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.IdleTimeout != d
				c.opts.IdleTimeout = d
			}

		case c.pClient + "/idle-check-frequency":
			if d, err := time.ParseDuration(v); err == nil {
				needUpdate = needUpdate || c.opts.IdleCheckFrequency != d
				c.opts.IdleCheckFrequency = d
			}
		}
	}

	return needUpdate
}

func (c *Connector) reset(last map[string]string) bool {
	needUpdate := false

	for k := range last {
		if strings.HasPrefix(k, c.pHost) {
			needUpdate = true
			c.opts.SentinelAddrs = []string{}

			break
		}
	}

	for k := range last {
		switch k {
		case c.pClient + "/master-name":
			needUpdate = true
			c.opts.MasterName = ""

		case c.pClient + "/password":
			needUpdate = true
			c.opts.Password = ""

		case c.pClient + "/db":
			needUpdate = true
			c.opts.DB = 0

		case c.pClient + "/max-retries":
			needUpdate = true
			c.opts.MaxRetries = 0

		case c.pClient + "/min-retry-backoff":
			needUpdate = true
			c.opts.MinRetryBackoff = 0

		case c.pClient + "/max-retry-backoff":
			needUpdate = true
			c.opts.MaxRetryBackoff = 0

		case c.pClient + "/dial-timeout":
			needUpdate = true
			c.opts.DialTimeout = 0

		case c.pClient + "/read-timeout":
			needUpdate = true
			c.opts.ReadTimeout = 0

		case c.pClient + "/write-timeout":
			needUpdate = true
			c.opts.WriteTimeout = 0

		case c.pClient + "/pool-size":
			needUpdate = true
			c.opts.PoolSize = 0

		case c.pClient + "/min-idle-conns":
			needUpdate = true
			c.opts.MinIdleConns = 0

		case c.pClient + "/max-conn-age":
			needUpdate = true
			c.opts.MaxConnAge = 0

		case c.pClient + "/pool-timeout":
			needUpdate = true
			c.opts.PoolTimeout = 0

		case c.pClient + "/idle-timeout":
			needUpdate = true
			c.opts.IdleTimeout = 0

		case c.pClient + "/idle-check-frequency":
			needUpdate = true
			c.opts.IdleCheckFrequency = 0
		}
	}

	return needUpdate
}

func clientConfig() []string {
	return []string{
		"master-name",
		"password",
		"db",
		"max-retries",
		"min-retry-backoff",
		"max-retry-backoff",
		"dial-timeout",
		"read-timeout",
		"write-timeout",
		"pool-size",
		"min-idle-conns",
		"max-conn-age",
		"pool-timeout",
		"idle-timeout",
		"idle-check-frequency",
	}
}
