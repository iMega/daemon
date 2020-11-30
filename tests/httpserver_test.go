package acceptance

import (
	"bytes"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/imega/daemon/tests/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reconnect http-server", func() {
	Context("current port 8080", func() {
		It("change to 8081", func() {
			expected := []byte("ok")
			actual := []byte("")

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "myhttp/http-server/host",
				Value: []byte("0.0.0.0:8081"),
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			for attempts := 30; attempts > 0; attempts-- {
				b, _ := helper.Request("http://appconsul:8081")

				if bytes.Compare(b, expected) == 0 {
					actual = b
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).To(Equal(expected))
		})

		It("change to 8080", func() {
			expected := []byte("ok")
			actual := []byte("")

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "myhttp/http-server/host",
				Value: []byte("0.0.0.0:8080"),
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			for attempts := 30; attempts > 0; attempts-- {
				b, _ := helper.Request("http://appconsul:8080")

				if bytes.Compare(b, expected) == 0 {
					actual = b
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).To(Equal(expected))
		})
	})
})
