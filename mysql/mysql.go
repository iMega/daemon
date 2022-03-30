package mysql

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/imega/daemon"
	"github.com/sirupsen/logrus"
)

// Connector is wrapper sql.DB.
type Connector struct {
	log logrus.FieldLogger

	pxHost   string
	pxClient string
	opts     *mysql.Config
	DB
	connMaxIdleTime time.Duration
	connMaxLifetime time.Duration
	maxIdleConns    int
	maxOpenConns    int

	WatcherConfigFuncs []daemon.WatcherConfigFunc
	daemon.HealthCheckFunc
	daemon.ShutdownFunc
}

const maxIdleConns = 2

// New get a instance of mysql.
func New(prefixHost, prefixClient string, l logrus.FieldLogger) *Connector {
	conn := &Connector{
		log: l,

		pxHost:   prefixHost + "/mysql/host",
		pxClient: prefixClient + "/mysql",
		opts:     mysql.NewConfig(),
		DB:       &fakerDB{},

		maxIdleConns: maxIdleConns,
	}

	conn.WatcherConfigFuncs = []daemon.WatcherConfigFunc{
		daemon.WatcherConfigFunc(func() daemon.WatcherConfig {
			return daemon.WatcherConfig{
				Prefix:    prefixHost,
				MainKey:   "mysql",
				Keys:      []string{"host"},
				ApplyFunc: conn.connect,
			}
		}),
		daemon.WatcherConfigFunc(func() daemon.WatcherConfig {
			return daemon.WatcherConfig{
				Prefix:    prefixClient,
				MainKey:   "mysql",
				Keys:      clientConfig(),
				ApplyFunc: conn.connect,
			}
		}),
	}

	conn.HealthCheckFunc = func() bool {
		if conn.DB == nil {
			conn.log.Error("failed to ping mysql")

			return false
		}

		if err := conn.DB.Ping(); err != nil {
			conn.log.Errorf("failed to ping mysql, %s", err)

			return false
		}

		conn.log.Debug("mysql ping is ok")

		return true
	}

	conn.ShutdownFunc = func() {
		if conn.DB == nil {
			conn.log.Error("failed to close connection to mysql")

			return
		}

		if err := conn.DB.Close(); err != nil {
			conn.log.Errorf("failed to close connection to mysql, %s", err)

			return
		}
	}

	return conn
}

func (db *Connector) connect(conf, last map[string]string) {
	reset := db.reset(last)
	config := db.config(conf)

	if !reset && !config {
		db.log.Debug("mysql connector has same configuration")

		return
	}

	conn, err := mysql.NewConnector(db.opts)
	if err != nil {
		db.log.Error(err)
	}

	if _, ok := db.DB.(*fakerDB); !ok {
		if err := db.DB.Close(); err != nil {
			db.log.Error(err)
		}

		db.log.Debug("mysql connection closed")

		db.DB = &fakerDB{}
	}

	db.DB = sql.OpenDB(conn)

	db.DB.SetMaxOpenConns(db.maxOpenConns)
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	db.DB.SetMaxIdleConns(db.maxIdleConns)
	db.DB.SetConnMaxIdleTime(db.connMaxIdleTime)

	db.log.Debug("mysql connection open")
}

func clientConfig() []string {
	return []string{
		"user",
		"password",
		"net",
		"db-name",
		"collation",
		"loc",
		"max-allowed-packet",
		"server-pub-key",
		"tls-config",
		"timeout",
		"read-timeout",
		"write-timeout",
		"allow-all-files",
		"allow-cleartext-passwords",
		"allow-native-passwords",
		"allow-old-passwords",
		"check-conn-liveness",
		"client-found-rows",
		"columns-with-alias",
		"interpolate-params",
		"multi-statements",
		"parse-time",
		"reject-read-only",
		"params",
		"conn-max-idle-time",
		"conn-max-lifetime",
		"max-idle-conns",
		"max-open-conns",
	}
}

const maxAllowedPacket = 4 << 20 // 4 MiB

func (db *Connector) reset(last map[string]string) bool {
	needUpdate := false

	for k := range last {
		switch k {
		case db.pxHost:
			needUpdate = true
			db.opts.Addr = "127.0.0.1:3306"

		case db.pxClient + "/user":
			needUpdate = true
			db.opts.User = ""

		case db.pxClient + "/password":
			needUpdate = true
			db.opts.Passwd = ""

		case db.pxClient + "/net":
			needUpdate = true
			db.opts.Net = "tcp"

		case db.pxClient + "/db-name":
			needUpdate = true
			db.opts.DBName = ""

		case db.pxClient + "/collation":
			needUpdate = true
			db.opts.Collation = "utf8mb4_general_ci"

		case db.pxClient + "/loc":
			needUpdate = true
			db.opts.Loc = time.UTC

		case db.pxClient + "/max-allowed-packet":
			needUpdate = true
			db.opts.MaxAllowedPacket = maxAllowedPacket

		case db.pxClient + "/server-pub-key":
			needUpdate = true
			db.opts.ServerPubKey = ""

		case db.pxClient + "/tls-config":
			needUpdate = true
			db.opts.TLSConfig = ""

		case db.pxClient + "/timeout":
			needUpdate = true
			db.opts.Timeout = 0

		case db.pxClient + "/read-timeout":
			needUpdate = true
			db.opts.ReadTimeout = 0

		case db.pxClient + "/write-timeout":
			needUpdate = true
			db.opts.WriteTimeout = 0

		case db.pxClient + "/allow-all-files":
			needUpdate = true
			db.opts.AllowAllFiles = false

		case db.pxClient + "/allow-cleartext-passwords":
			needUpdate = true
			db.opts.AllowCleartextPasswords = false

		case db.pxClient + "/allow-native-passwords":
			needUpdate = true
			db.opts.AllowNativePasswords = true

		case db.pxClient + "/allow-old-passwords":
			needUpdate = true
			db.opts.AllowOldPasswords = false

		case db.pxClient + "/check-conn-liveness":
			needUpdate = true
			db.opts.CheckConnLiveness = true

		case db.pxClient + "/client-found-rows":
			needUpdate = true
			db.opts.ClientFoundRows = false

		case db.pxClient + "/columns-with-alias":
			needUpdate = true
			db.opts.ColumnsWithAlias = false

		case db.pxClient + "/interpolate-params":
			needUpdate = true
			db.opts.InterpolateParams = false

		case db.pxClient + "/multi-statements":
			needUpdate = true
			db.opts.MultiStatements = false

		case db.pxClient + "/parse-time":
			needUpdate = true
			db.opts.ParseTime = false

		case db.pxClient + "/reject-read-only":
			needUpdate = true
			db.opts.RejectReadOnly = false

		case db.pxClient + "/params":
			needUpdate = true
			db.opts.Params = make(map[string]string)

		case db.pxClient + "/conn-max-idle-time":
			needUpdate = true
			db.connMaxIdleTime = 0

		case db.pxClient + "/conn-max-lifetime":
			needUpdate = true
			db.connMaxLifetime = 0

		case db.pxClient + "/max-idle-conns":
			needUpdate = true
			db.maxIdleConns = 0

		case db.pxClient + "/max-open-conns":
			needUpdate = true
			db.maxOpenConns = 0
		}
	}

	return needUpdate
}

const valueTrue = "true"

func (db *Connector) config(conf map[string]string) bool {
	needUpdate := false

	for key, value := range conf {
		switch key {
		case db.pxHost:
			needUpdate = needUpdate || db.opts.Addr != value
			db.opts.Addr = value

		case db.pxClient + "/user":
			needUpdate = needUpdate || db.opts.User != value
			db.opts.User = value

		case db.pxClient + "/password":
			needUpdate = needUpdate || db.opts.Passwd != value
			db.opts.Passwd = value

		case db.pxClient + "/net":
			needUpdate = needUpdate || db.opts.Net != value
			db.opts.Net = value

		case db.pxClient + "/db-name":
			needUpdate = needUpdate || db.opts.DBName != value
			db.opts.DBName = value

		case db.pxClient + "/collation":
			needUpdate = needUpdate || db.opts.Collation != value
			db.opts.Collation = value

		case db.pxClient + "/loc":
			if loc, err := time.LoadLocation(value); err == nil {
				needUpdate = needUpdate || db.opts.Loc != loc
				db.opts.Loc = loc
			}

		case db.pxClient + "/max-allowed-packet":
			if i, err := strconv.Atoi(value); err == nil {
				needUpdate = needUpdate || db.opts.MaxAllowedPacket != i
				db.opts.MaxAllowedPacket = i
			}

		case db.pxClient + "/server-pub-key":
			needUpdate = needUpdate || db.opts.ServerPubKey != value
			db.opts.ServerPubKey = value

		case db.pxClient + "/tls-config":
			needUpdate = needUpdate || db.opts.TLSConfig != value
			db.opts.TLSConfig = value

		case db.pxClient + "/timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || db.opts.Timeout != d
				db.opts.Timeout = d
			}

		case db.pxClient + "/read-timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || db.opts.ReadTimeout != d
				db.opts.ReadTimeout = d
			}

		case db.pxClient + "/write-timeout":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || db.opts.WriteTimeout != d
				db.opts.WriteTimeout = d
			}

		case db.pxClient + "/allow-all-files":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.AllowAllFiles != val
			db.opts.AllowAllFiles = value == valueTrue

		case db.pxClient + "/allow-cleartext-passwords":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.AllowCleartextPasswords != val
			db.opts.AllowCleartextPasswords = value == valueTrue

		case db.pxClient + "/allow-native-passwords":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.AllowNativePasswords != val
			db.opts.AllowNativePasswords = value == valueTrue

		case db.pxClient + "/allow-old-passwords":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.AllowOldPasswords != val
			db.opts.AllowOldPasswords = value == valueTrue

		case db.pxClient + "/check-conn-liveness":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.CheckConnLiveness != val
			db.opts.CheckConnLiveness = value == valueTrue

		case db.pxClient + "/client-found-rows":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.ClientFoundRows != val
			db.opts.ClientFoundRows = value == valueTrue

		case db.pxClient + "/columns-with-alias":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.ColumnsWithAlias != val
			db.opts.ColumnsWithAlias = value == valueTrue

		case db.pxClient + "/interpolate-params":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.InterpolateParams != val
			db.opts.InterpolateParams = value == valueTrue

		case db.pxClient + "/multi-statements":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.MultiStatements != val
			db.opts.MultiStatements = value == valueTrue

		case db.pxClient + "/parse-time":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.ParseTime != val
			db.opts.ParseTime = val

		case db.pxClient + "/reject-read-only":
			val := value == valueTrue
			needUpdate = needUpdate || db.opts.RejectReadOnly != val
			db.opts.RejectReadOnly = value == valueTrue

		case db.pxClient + "/params":
			m := map[string]string{}
			if err := json.Unmarshal([]byte(value), &m); err == nil {
				needUpdate = needUpdate || !reflect.DeepEqual(db.opts.Params, m)
				db.opts.Params = m
			}

		case db.pxClient + "/conn-max-idle-time":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || db.connMaxIdleTime != d
				db.connMaxIdleTime = d
			}

		case db.pxClient + "/conn-max-lifetime":
			if d, err := time.ParseDuration(value); err == nil {
				needUpdate = needUpdate || db.connMaxLifetime != d
				db.connMaxLifetime = d
			}

		case db.pxClient + "/max-idle-conns":
			if i, err := strconv.Atoi(value); err == nil {
				needUpdate = needUpdate || db.maxIdleConns != i
				db.maxIdleConns = i
			}

		case db.pxClient + "/max-open-conns":
			if i, err := strconv.Atoi(value); err == nil {
				needUpdate = needUpdate || db.maxOpenConns != i
				db.maxOpenConns = i
			}
		}
	}

	return needUpdate
}
