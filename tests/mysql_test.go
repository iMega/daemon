package acceptance

import (
	"bytes"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/imega/daemon/tests/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reconnect mysql", func() {
	Context("current instance 1", func() {
		It("change instance 2", func() {
			expected := []byte("mysql2")
			actual := []byte("")

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "instance/mysql/host",
				Value: expected,
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			for attempts := 30; attempts > 0; attempts-- {

				b, _ := helper.Request("http://appconsul:80/test_mysql_reconnect_between_instances")

				if bytes.Compare(b, expected) == 0 {
					actual = b
					break
				}

				<-time.After(1 * time.Second)
			}

			Expect(actual).To(Equal(expected))
		})

		It("change instance 1", func() {
			expected := []byte("mysql1")
			actual := []byte("")

			client, err := api.NewClient(api.DefaultConfig())
			Expect(err).NotTo(HaveOccurred())

			_, err = client.KV().Put(&api.KVPair{
				Key:   "instance/mysql/host",
				Value: expected,
			}, &api.WriteOptions{})
			Expect(err).NotTo(HaveOccurred())

			for attempts := 30; attempts > 0; attempts-- {

				b, _ := helper.Request("http://appconsul:80/test_mysql_reconnect_between_instances")

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
