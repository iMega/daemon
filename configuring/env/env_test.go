package env

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/imega/daemon"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	type args struct {
		key string
	}

	filename := "example-env."
	content := "value"

	os.Setenv("TEST_ENV", content)
	tmpfile, err := ioutil.TempFile(os.TempDir(), filename)
	os.Setenv("TEST_ENV_2_FILE", tmpfile.Name())

	if err != nil {
		t.Errorf("failed to create temp file, %s", err)
	}

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Errorf("failed to write to temp file, %s", err)
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "read env",
			args: args{
				key: "TEST_ENV",
			},
			want: content,
		},
		{
			name: "read env from file",
			args: args{
				key: "TEST_ENV_2",
			},
			want: content,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}

	os.Remove(tmpfile.Name())
}

func Test_watcher_Read(t *testing.T) {
	type fields struct {
		f []daemon.WatcherConfigFunc
	}
	var actual map[string]string

	os.Setenv("MY_DAEMON_HTTP_SERVER_READ_HEADER_TIMEOUT", "20")
	os.Setenv("MY_DAEMON_HTTP_SERVER_WRITE_TIMEOUT", "10")

	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				f: []daemon.WatcherConfigFunc{
					func() daemon.WatcherConfig {
						return daemon.WatcherConfig{
							Prefix:  "my-daemon",
							MainKey: "http-server",
							Keys: []string{
								"read-header-timeout",
								"write-timeout",
							},
							ApplyFunc: func(c, r map[string]string) {
								actual = c
							},
						}
					},
				},
			},
			want: map[string]string{
				"my-daemon/http-server/read-header-timeout": "20",
				"my-daemon/http-server/write-timeout":       "10",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Watcher{
				f: tt.fields.f,
			}
			if err := w.Read(); (err != nil) != tt.wantErr {
				t.Errorf("watcher.Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !assert.Equal(t, tt.want, actual) {
				t.Error("watcher.Read() failed assertion")
			}
		})
	}
}
