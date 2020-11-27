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

package consul

import (
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

func newConsulLogger(log logrus.FieldLogger) hclog.Logger {
	l := &logger{log: log}

	l.hlog = l

	return l
}

type logger struct {
	log  logrus.FieldLogger
	hlog hclog.Logger
}

func (l *logger) Log(level hclog.Level, msg string, args ...interface{}) {
	if e, ok := l.log.(*logrus.Entry); ok {
		e.Log(logrus.Level(level), msg, args)
	}
}

func (l *logger) Trace(msg string, args ...interface{}) {
	if e, ok := l.log.(*logrus.Entry); ok {
		e.Trace(msg, args)
	}
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.log.Debug(msg, args)
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.log.Info(msg, args)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.log.Warn(msg, args)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.log.Error(msg, args)
}

func (l *logger) IsTrace() bool {
	return l.log.(*logrus.Entry).Logger.Level <= logrus.DebugLevel
}

func (l *logger) IsDebug() bool {
	if e, ok := l.log.(*logrus.Entry); ok {
		return e.Level <= logrus.DebugLevel
	}

	return false
}

func (l *logger) IsInfo() bool {
	if e, ok := l.log.(*logrus.Entry); ok {
		return e.Level <= logrus.InfoLevel
	}

	return false
}

func (l *logger) IsWarn() bool {
	if e, ok := l.log.(*logrus.Entry); ok {
		return e.Level <= logrus.WarnLevel
	}

	return false
}

func (l *logger) IsError() bool {
	if e, ok := l.log.(*logrus.Entry); ok {
		return e.Level <= logrus.ErrorLevel
	}

	return false
}

func (l *logger) ImpliedArgs() []interface{} {
	return nil
}

func (l *logger) With(args ...interface{}) hclog.Logger {
	return nil
}

func (l *logger) Name() string {
	return ""
}

func (l *logger) Named(name string) hclog.Logger {
	return l.hlog
}

func (l *logger) ResetNamed(name string) hclog.Logger {
	return l.hlog
}

func (l *logger) SetLevel(level hclog.Level) {}

func (l *logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return nil
}

func (l *logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	if e, ok := l.log.(*logrus.Entry); ok {
		return e.Logger.Out
	}

	return os.Stdout
}
