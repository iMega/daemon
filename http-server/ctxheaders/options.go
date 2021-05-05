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

package ctxheaders

var defaultOptions = &options{
	HeadersToTags:    nil,
	HeadersToContext: map[string]string{},
}

type options struct {
	HeadersToTags    []string
	HeadersToContext map[string]string
}

type Option func(*options)

func HeadersToTags(v []string) Option {
	return func(o *options) {
		o.HeadersToTags = v
	}
}

// HeadersToContext is a map relationships between headers and context.
func HeadersToContext(v map[string]string) Option {
	return func(o *options) {
		o.HeadersToContext = v
	}
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions

	for _, o := range opts {
		o(optCopy)
	}

	return optCopy
}
