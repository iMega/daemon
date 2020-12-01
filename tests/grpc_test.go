package acceptance

import (
	"context"
	"time"

	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var _ = Describe("reconnect grpc-server", func() {
	Context("current port 9001", func() {
		It("change to 9002", func() {
			actual := grpc_health_v1.HealthCheckResponse_UNKNOWN

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "testclient/grpc/host",
				Value: []byte("0.0.0.0:9002"),
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			cc, err := grpc.Dial("appconsul:9002", grpc.WithInsecure())
			Expect(err).NotTo(HaveOccurred())

			hc := grpc_health_v1.NewHealthClient(cc)
			for attempts := 30; attempts > 0; attempts-- {

				resp, _ := hc.Check(
					context.Background(),
					&grpc_health_v1.HealthCheckRequest{},
				)

				if resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
					actual = resp.GetStatus()
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).To(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
		})

		It("change to 9001", func() {
			actual := grpc_health_v1.HealthCheckResponse_UNKNOWN

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "testclient/grpc/host",
				Value: []byte("0.0.0.0:9001"),
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			cc, err := grpc.Dial("appconsul:9001", grpc.WithInsecure())
			Expect(err).NotTo(HaveOccurred())

			hc := grpc_health_v1.NewHealthClient(cc)
			for attempts := 30; attempts > 0; attempts-- {

				resp, _ := hc.Check(
					context.Background(),
					&grpc_health_v1.HealthCheckRequest{},
				)

				if resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
					actual = resp.GetStatus()
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).To(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
		})
	})
})
