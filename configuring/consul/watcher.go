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
	"fmt"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/imega/daemon"
	"github.com/sirupsen/logrus"
)

type Watcher struct {
	log           logrus.FieldLogger
	wathFunc      []daemon.WatcherConfigFunc
	LastConfMutex sync.RWMutex
	LastConf      map[string]map[string]string
}

// Watch .
func Watch(log logrus.FieldLogger, f ...daemon.WatcherConfigFunc) *Watcher {
	return &Watcher{
		log:           log,
		wathFunc:      f,
		LastConfMutex: sync.RWMutex{},
		LastConf:      make(map[string]map[string]string),
	}
}

func (w *Watcher) Read() error {
	hlog := newConsulLogger(w.log)
	conf := api.DefaultConfigWithLogger(hlog)

	for _, fn := range w.wathFunc {
		wConf := fn()

		prefixKey := wConf.Prefix + "/" + wConf.MainKey

		plan, err := watch.Parse(map[string]interface{}{
			"type":   "keyprefix",
			"prefix": prefixKey,
		})
		if err != nil {
			return fmt.Errorf("failed to parse keys: %w", err)
		}

		plan.Logger = hlog

		func(p *watch.Plan, cbFunc daemon.ApplyConfigFunc, key string) {
			p.HybridHandler = func(v watch.BlockingParamVal, i interface{}) {
				conf := make(map[string]string)

				if pairs, ok := i.(api.KVPairs); ok {
					for _, pair := range pairs {
						conf[pair.Key] = string(pair.Value)
					}
				}

				if len(conf) == 0 {
					return
				}

				w.LastConfMutex.RLock()
				cbFunc(conf, keys4reset(conf, w.LastConf[key]))
				w.LastConfMutex.RUnlock()

				w.LastConfMutex.Lock()
				w.LastConf[key] = conf
				w.LastConfMutex.Unlock()
			}
		}(plan, wConf.ApplyFunc, prefixKey)

		go func() {
			if err := plan.RunWithConfig(conf.Address, conf); err != nil {
				w.log.Error(err)
			}
		}()
	}

	return nil
}

func keys4reset(current, last map[string]string) map[string]string {
	reset := make(map[string]string)

	for k, v := range last {
		if _, ok := current[k]; !ok {
			reset[k] = v
		}
	}

	return reset
}
