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

package health

import (
	"github.com/imega/daemon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// New registers a health-check server and its implementation
// to the gRPC server. This must be called before invoking Serve.
func New(s *grpc.Server, f ...daemon.HealthCheckFunc) {
	grpc_health_v1.RegisterHealthServer(s, &server{f})
	reflection.Register(s)
}

type server struct {
	fn []daemon.HealthCheckFunc
}

func (s *server) Check(
	context.Context,
	*grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	status := grpc_health_v1.HealthCheckResponse_SERVING

	for _, p := range s.fn {
		if ok := p(); !ok {
			status = grpc_health_v1.HealthCheckResponse_NOT_SERVING

			break
		}
	}

	return &grpc_health_v1.HealthCheckResponse{Status: status}, nil
}

func (s *server) Watch(
	*grpc_health_v1.HealthCheckRequest,
	grpc_health_v1.Health_WatchServer,
) error {
	return nil
}
