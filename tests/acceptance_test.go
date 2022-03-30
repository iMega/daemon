package acceptance

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var _ = BeforeSuite(func() {
	var errConsul, errEnv error
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		errConsul = WaitingForSystemUnderTestBeReady("appconsul:9000")
	}()

	go func() {
		defer wg.Done()
		errEnv = WaitingForSystemUnderTestBeReady("appenv:9000")
	}()

	wg.Wait()

	Expect(errConsul).NotTo(HaveOccurred())
	Expect(errEnv).NotTo(HaveOccurred())
})

func WaitingForSystemUnderTestBeReady(host string) error {
	cc, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return err
	}
	hc := grpc_health_v1.NewHealthClient(cc)
	for attempts := 30; attempts > 0; attempts-- {
		resp, err := hc.Check(
			context.Background(),
			&grpc_health_v1.HealthCheckRequest{},
		)
		if err == nil && resp != nil && resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
			return nil
		}
		log.Printf("ATTEMPTING TO CONNECT")
		<-time.After(1 * time.Second)
	}

	return errors.New("SUT is not ready for tests")
}

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}
