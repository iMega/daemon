package acceptance

import (
	"bytes"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/imega/daemon/tests/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reconnect redis", func() {
	Context("current instance 0", func() {
		It("change to instance 1", func() {
			expected := []byte(":16379")
			actual := []byte("")

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "instance/redis-sentinel/host/instance-1",
				Value: []byte("sentinel1:26379"),
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			for attempts := 30; attempts > 0; attempts-- {
				b, _ := helper.Request("http://appconsul:80/test_redis_reconnect_between_instances")

				if bytes.Contains(b, expected) {
					actual = b
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).Should(ContainSubstring(string(expected)))
		})

		It("change to instance 0", func() {
			expected := []byte(":6379")
			actual := []byte("")

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "instance/redis-sentinel/host/instance-1",
				Value: []byte("sentinel0:26379"),
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			for attempts := 30; attempts > 0; attempts-- {
				b, _ := helper.Request("http://appconsul:80/test_redis_reconnect_between_instances")

				if bytes.Contains(b, expected) {
					actual = b
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).Should(ContainSubstring(string(expected)))
		})
	})
})
