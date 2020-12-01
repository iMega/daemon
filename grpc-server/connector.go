package grpcserver

import (
	"errors"
	"fmt"
	"net"
	"runtime/debug"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/imega/daemon"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Connector is a wrapped grpc server.
type Connector struct {
	prefixClient string

	log   logrus.FieldLogger
	opts  *optionsServer
	gOpts []grpc.ServerOption

	srv *grpc.Server
	rs  RegisterServices
	lis net.Listener

	daemon.WatcherConfigFunc
	daemon.ShutdownFunc
}

type optionsServer struct {
	Host string
}

const (
	configGRPCHost = "/grpc/host"
)

// RegisterServices will register any services.
type RegisterServices func(srv *grpc.Server)

// Option .
type Option func(*Connector)

// WithLogger .
func WithLogger(l logrus.FieldLogger) Option {
	return func(o *Connector) {
		if e, ok := l.(*logrus.Entry); ok {
			newVerbosityLogger(e)
		}

		o.log = l
	}
}

// WithServices .
func WithServices(rs RegisterServices) Option {
	return func(o *Connector) {
		o.rs = rs
	}
}

// WithServerOptions .
func WithServerOptions(opts ...grpc.ServerOption) Option {
	return func(o *Connector) {
		o.gOpts = opts
	}
}

// New get a instance of grpc.
func New(prefix string, opts ...Option) *Connector {
	s := &Connector{
		prefixClient: prefix,

		opts: &optionsServer{
			Host: "0.0.0.0:65535",
		},
	}

	for _, o := range opts {
		o(s)
	}

	s.srv = s.newServer()

	s.WatcherConfigFunc = func() daemon.WatcherConfig {
		return daemon.WatcherConfig{
			Prefix:    prefix,
			MainKey:   "grpc",
			Keys:      []string{"host"},
			ApplyFunc: s.connect,
		}
	}

	s.ShutdownFunc = func() {
		s.srv.GracefulStop()
	}

	return s
}

var errRecovery = errors.New("recovery handler error")

func (s *Connector) newServer() *grpc.Server {
	var log *logrus.Entry
	if e, ok := s.log.(*logrus.Entry); ok {
		log = e
	}

	rOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			stack := string(debug.Stack())

			return fmt.Errorf("%w: %s: %s", errRecovery, p, stack)
		}),
	}

	loggerOpts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(methodFullName string, err error) bool {
			if err == nil && methodFullName == "/grpc.health.v1.Health/Check" {
				return false
			}

			return true
		}),
	}

	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(log, loggerOpts...),
			grpc_recovery.UnaryServerInterceptor(rOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(rOpts...),
		),
	}

	return grpc.NewServer(append(opts, s.gOpts...)...)
}

func (s *Connector) connect(conf, last map[string]string) {
	reset := s.reset(last)
	config := s.config(conf)

	if !reset && !config {
		s.log.Debug("grpc connector has same configuration")

		return
	}

	if s.lis != nil {
		s.log.Debugf("grpc connector start graceful stop, %s", s.lis.Addr().String())
		s.srv.GracefulStop()
		s.log.Debug("grpc connector end graceful stop")

		s.srv = s.newServer()
	}

	l, err := net.Listen("tcp", s.opts.Host)
	if err != nil {
		s.log.Errorf(
			"failed to listen on the TCP network address %s, %s",
			s.opts.Host,
			err,
		)

		return
	}

	s.lis = l

	if s.rs != nil {
		s.rs(s.srv)
	}

	go func() {
		s.log.Debugf("grpc connector start on %s", l.Addr().String())

		if err := s.srv.Serve(l); err != nil {
			if !errors.Is(err, grpc.ErrServerStopped) {
				s.log.Error(err)
			}
		}
	}()
}

func (s *Connector) config(conf map[string]string) bool {
	needUpdate := false

	for k, v := range conf {
		if k == s.prefixClient+configGRPCHost {
			needUpdate = needUpdate || s.opts.Host != v
			s.opts.Host = v
		}
	}

	return needUpdate
}

func (s *Connector) reset(last map[string]string) bool {
	needUpdate := false

	for k := range last {
		if k == s.prefixClient+configGRPCHost {
			needUpdate = true
			s.opts.Host = "0.0.0.0:65535"
		}
	}

	return needUpdate
}
