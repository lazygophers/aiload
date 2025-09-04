package memqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_apply(t *testing.T) {
	t.Run("default_values", func(t *testing.T) {
		cfg := &Config[any]{}
		cfg.apply()

		assert.Equal(t, 500, cfg.MaxSize)
		assert.Equal(t, 1, cfg.ConcurrentCount)
		assert.Equal(t, 3, cfg.MaxRetry)
	})

	t.Run("custom_values", func(t *testing.T) {
		cfg := &Config[any]{
			MaxSize:         1000,
			ConcurrentCount: 5,
			MaxRetry:        10,
			Delay:           time.Second,
		}
		cfg.apply()

		assert.Equal(t, 1000, cfg.MaxSize)
		assert.Equal(t, 5, cfg.ConcurrentCount)
		assert.Equal(t, 10, cfg.MaxRetry)
	})

	t.Run("zero_values", func(t *testing.T) {
		cfg := &Config[any]{
			MaxSize:         0,
			ConcurrentCount: 0,
			MaxRetry:        0,
		}
		cfg.apply()

		assert.Equal(t, 500, cfg.MaxSize)
		assert.Equal(t, 1, cfg.ConcurrentCount)
		assert.Equal(t, 3, cfg.MaxRetry)
	})
}
