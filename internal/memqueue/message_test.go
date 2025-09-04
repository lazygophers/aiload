package memqueue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage_init(t *testing.T) {
	t.Run("default_values", func(t *testing.T) {
		cfg := &Config[any]{
			Delay: time.Second,
		}
		msg := &Message[any]{}
		msg.init(cfg)

		assert.NotZero(t, msg.CreatedAt)
		assert.NotEmpty(t, msg.Id)
		assert.NotZero(t, msg.ExecAt)
		assert.Equal(t, 0, msg.RetryCount)
	})

	t.Run("with_custom_values", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		execAt := time.Now().Add(time.Hour)
		id := "custom-id"

		cfg := &Config[any]{}
		msg := &Message[any]{
			CreatedAt: createdAt,
			Id:        id,
			ExecAt:    execAt,
		}
		msg.init(cfg)

		assert.Equal(t, createdAt, msg.CreatedAt)
		assert.Equal(t, id, msg.Id)
		assert.Equal(t, execAt, msg.ExecAt)
	})

	t.Run("zero_delay", func(t *testing.T) {
		cfg := &Config[any]{
			Delay: 0,
		}
		msg := &Message[any]{}
		msg.init(cfg)

		assert.True(t, msg.ExecAt.IsZero())
	})
}
