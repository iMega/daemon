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

package logging

import "context"

// Logger is an interface for logger.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// Config is a configuration logger.
type Config struct {
	Channel string
	BuildID string
	Level   string
}

func ContextWithLogger(ctx context.Context, log Logger) context.Context {
	if _, ok := ctx.Value(key{}).(Logger); ok {
		return ctx
	}

	return context.WithValue(ctx, key{}, log)
}

func GetLogger(ctx context.Context) Logger {
	val, ok := ctx.Value(key{}).(Logger)
	if !ok {
		return &noopLog{}
	}

	return val
}

type key struct{}

func GetNoopLog() Logger {
	return &noopLog{}
}

type noopLog struct{}

func (nl *noopLog) Infof(string, ...interface{})  {}
func (nl *noopLog) Errorf(string, ...interface{}) {}
func (nl *noopLog) Debugf(string, ...interface{}) {}
