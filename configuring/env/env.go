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

package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/imega/daemon"
)

// Read retrieves the value of the environment variable named by the key.
// As an alternative to passing sensitive information via environment variables,
// _FILE may be appended to the previously listed environment variables,
// causing the initialization script to load the values for those variables
// from files present in the container. In particular, this can be used
// to load passwords from Docker secrets stored in /run/secrets/<secret_name>
// files.
//
// https://docs.docker.com/engine/swarm/secrets/
// https://docs.docker.com/compose/compose-file/compose-versioning/#version-31
//
// docker-compose.yml
// ------------------
// version: "3.1"
//
// secrets:
//   mypassword:
//     external: true
//
// services:
//   app:
//     image: your-app-in-docker-image
//     environment:
//       - PASSWORD_FILE=/run/secrets/mypassword
//       - LOGIN=mylogin
//
// EXAMPLE:
//
// passwd, _ := env.Read("PASSWORD") // reading from PASSWORD_FILE
// login, _ := end.Read("LOGIN")     // reading from LOGIN.
func Read(key string) (string, error) {
	if filename := os.Getenv(key + "_FILE"); filename != "" {
		value, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		return strings.TrimSpace(string(value)), nil
	}

	return os.Getenv(key), nil
}

type watcher struct {
	f []daemon.WatcherConfigFunc
}

// Once .
func Once(f ...daemon.WatcherConfigFunc) daemon.ConfigReader {
	return &watcher{f}
}

func (w *watcher) Read() error {
	for _, fn := range w.f {
		mapKeys := make(map[string]string)
		wConf := fn()

		for _, k := range wConf.Keys {
			env := wConf.Prefix + "_" + wConf.MainKey + "_" + k
			env = strings.ReplaceAll(env, "-", "_")
			v, _ := Read(env)
			if v != "" {
				mapKeys[wConf.Prefix+k] = v
			}
		}

		wConf.ApplyFunc(mapKeys, nil)
	}

	return nil
}
