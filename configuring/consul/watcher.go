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

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/imega/daemon"
	"github.com/sirupsen/logrus"
)

type watcher struct {
	log      logrus.FieldLogger
	wathFunc []daemon.WatcherConfigFunc
	LastConf map[string]map[string]string
}

// Watch .
func Watch(log logrus.FieldLogger, f ...daemon.WatcherConfigFunc) daemon.ConfigReader {
	return &watcher{log: log, wathFunc: f}
}

func (w *watcher) Read() error {
	hlog := newConsulLogger(w.log)
	conf := api.DefaultConfigWithLogger(hlog)

	for _, fn := range w.wathFunc {
		keys, cb := fn()

		for _, k := range keys {
			plan, err := watch.Parse(map[string]interface{}{
				"type":   "keyprefix",
				"prefix": k,
			})
			if err != nil {
				return fmt.Errorf("failed to parse keys: %w", err)
			}

			plan.Logger = hlog

			func(p *watch.Plan, cb daemon.ApplyConfigFunc, k string) {
				p.HybridHandler = func(v watch.BlockingParamVal, i interface{}) {
					m := make(map[string]string)

					if pairs, ok := i.(api.KVPairs); ok {
						for _, pair := range pairs {
							m[pair.Key] = string(pair.Value)
						}
					}

					if len(m) == 0 {
						return
					}

					cb(m, keys4reset(m, w.LastConf[k]))

					w.LastConf[k] = m
				}
			}(plan, cb, k)

			go func() {
				if err := plan.RunWithConfig(conf.Address, conf); err != nil {
					w.log.Error(err)
				}
			}()
		}
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
