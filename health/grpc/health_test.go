package health

import (
	"testing"

	"github.com/imega/daemon"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestRegisterResourceHealthCheck_RegistersHealthCheck(t *testing.T) {
	s := server{
		fn: []daemon.HealthCheckFunc{
			func() bool { return true },
		},
	}

	if len(s.fn) != 1 {
		t.Fatalf("RegisterHealthCheckFunc does not adds health check funcs to server probes")
	}
}

func TestCheck_ExecutesResourceHealthCheckFuncs(t *testing.T) {
	executed := false
	s := server{
		fn: []daemon.HealthCheckFunc{
			func() bool {
				executed = true

				return executed
			},
		},
	}

	_, err := s.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("health check Check failed")
	}

	if executed != true {
		t.Fatalf("health check server does not executes registered health checks")
	}
}

func TestCheck_AllChecksPassed_ReportsStatusSERVING(t *testing.T) {
	s := server{
		fn: []daemon.HealthCheckFunc{
			func() bool { return true },
		},
	}

	resp, err := s.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("health check server does not expected to return error")
	}

	if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("health check expected to report SERVING status")
	}
}

func TestCheck_AnyCheckFailed_ReportsStatusNOT_SERVING(t *testing.T) {
	s := server{
		fn: []daemon.HealthCheckFunc{
			func() bool { return true },
			func() bool { return false },
		},
	}

	resp, _ := s.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})

	if resp.GetStatus() != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("health check expected to report NOT_SERVING status")
	}
}
