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

package daemon

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type daemon struct {
	Log logrus.FieldLogger

	sf []ShutdownFunc
}

// Daemon is a interface.
type Daemon interface {
	Run(shutdownTimeout time.Duration) error
	RegisterShutdownFunc(f ShutdownFunc)
}

// New create a new Daemon.
func New(l logrus.FieldLogger) Daemon {
	d := &daemon{
		Log: l,
	}

	return d
}

// Run daemon.
func (d *daemon) Run(shutdownTimeout time.Duration) error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)
	<-sigchan

	return d.shutdown(shutdownTimeout)
}

// ShutdownFunc .
type ShutdownFunc func()

// RegisterShutdownFunc .
func (d *daemon) RegisterShutdownFunc(f ShutdownFunc) {
	d.sf = append(d.sf, f)
}

// ErrShutdownTimeout .
var ErrShutdownTimeout = errors.New("shutdown timeout")

func (d *daemon) shutdown(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	wg := sync.WaitGroup{}

	for _, f := range d.sf {
		wg.Add(1)

		go func(f ShutdownFunc) {
			defer wg.Done()
			f()
		}(f)
	}

	doneChan := make(chan struct{})

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	select {
	case <-timer.C:
		return ErrShutdownTimeout
	case <-doneChan:
		return nil
	}
}
