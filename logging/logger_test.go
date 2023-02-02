package logging_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/imega/daemon/logging/wrapzerolog"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := wrapzerolog.New(zerolog.New(buf).With().Logger())

	logger.Infof("test")

	expected := map[string]interface{}{
		"level":   "info",
		"message": "test",
	}

	actual := map[string]interface{}{}

	err := json.NewDecoder(buf).Decode(&actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestWithFields(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := wrapzerolog.New(zerolog.New(buf).With().Logger())

	ctxLogger := logger.WithFields(map[string]interface{}{
		"tag": "test-tag",
	})

	ctxLogger.Infof("test")

	expected := map[string]interface{}{
		"level":   "info",
		"tag":     "test-tag",
		"message": "test",
	}

	actual := map[string]interface{}{}

	err := json.NewDecoder(buf).Decode(&actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
