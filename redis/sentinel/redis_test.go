package redis

import (
	"sort"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestConnector_config(t *testing.T) {
	type fields struct {
		opts    *redis.FailoverOptions
		pHost   string
		pClient string
	}
	type args struct {
		conf map[string]string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantOpts *redis.FailoverOptions
	}{
		{
			name: "set host",
			fields: fields{
				opts:    &redis.FailoverOptions{},
				pHost:   "instance",
				pClient: "alient",
			},
			args: args{
				conf: map[string]string{
					"instance/redis-sentinel/host/instance-1": "host-1",
					"instance/redis-sentinel/host/instance-2": "host-2",
					"instance/redis-sentinel/host/instance-3": "host-3",
				},
			},
			want: true,
			wantOpts: &redis.FailoverOptions{
				SentinelAddrs: []string{"host-1", "host-2", "host-3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connector{
				opts:    tt.fields.opts,
				pHost:   tt.fields.pHost,
				pClient: tt.fields.pClient,
			}
			if got := c.config(tt.args.conf); got != tt.want {
				t.Errorf("Connector.config() = %v, want %v", got, tt.want)
			}

			sort.Strings(c.opts.SentinelAddrs)
			assert.Equal(t, tt.wantOpts, c.opts)
		})
	}
}
