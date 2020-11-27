package env

import (
	"io/ioutil"
	"os"
	"testing"
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
