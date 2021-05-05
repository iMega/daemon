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

import (
	"context"
	"net/http"

	http_ctxtags "github.com/improbable-eng/go-httpwares/tags"
)

// Middleware is a server-side http ware for set headers in tags
// and/or in context.
//
// Example
//
// mux := http.NewServeMux()
// mux.HandleFunc("/", handler)
// httpServer := &http.Server{
//     Addr:    "0.0.0.0:8080",
//     Handler: http_ctxtags.Middleware("http")(
//         ctxheaders.Middleware(
//             mux,
//             ctxheaders.HeadersToContext(map[string]string{
//                 "X-Site-ID": "x-site-id",
//             }),
//         ),
//     ),
// }.
func Middleware(next http.Handler, opts ...Option) http.Handler {
	o := evaluateOptions(opts)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx  = r.Context()
			tags = http_ctxtags.ExtractInbound(r)
		)

		for k, v := range o.HeadersToContext {
			val := r.Header.Get(k)
			tags.Set(v, val)
			ctx = context.WithValue(ctx, v, val) // nolint:golint,staticcheck
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
