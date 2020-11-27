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

import "github.com/sirupsen/logrus"

// Config is a configuration logger.
type Config struct {
	Channel       string
	BuildID       string
	Level         string
	TextFormatter *logrus.TextFormatter
	JSONFormatter *logrus.JSONFormatter
}

// New create a new logger.
func New(conf Config) logrus.FieldLogger {
	if conf.Level == "" {
		conf.Level = "error"
	}

	logLevel, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)

	if conf.TextFormatter != nil {
		logrus.SetFormatter(conf.TextFormatter)
	}

	if conf.JSONFormatter != nil {
		logrus.SetFormatter(conf.JSONFormatter)
	}

	return logrus.WithFields(
		logrus.Fields{
			"channel":  conf.Channel,
			"build_id": conf.BuildID,
		},
	)
}
