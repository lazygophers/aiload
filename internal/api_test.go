package internal

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/bbolt"
)

func TestHealthCheck(t *testing.T) {
	// given
	app := fiber.New()
	state, err := NewState()
	assert.NoError(t, err)
	RegisterRoutes(app, state)

	// when
	resp, _ := app.Test(httptest.NewRequest("GET", "/health", nil))

	// then
	t.Cleanup(func() {
		assert.NoError(t, state.DB.Close())
		assert.NoError(t, state.Cache.Close())
		assert.NoError(t, os.RemoveAll("./data"))
	})

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(t, `{"status":"ok"}`, string(body))
}

func TestHelloWorld(t *testing.T) {
	// given
	app := fiber.New()
	state, err := NewState()
	assert.NoError(t, err)
	RegisterRoutes(app, state)

	// when
	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))

	// then
	t.Cleanup(func() {
		assert.NoError(t, state.DB.Close())
		assert.NoError(t, state.Cache.Close())
		assert.NoError(t, os.RemoveAll("./data"))
	})

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "Hello, World!", string(body))
}

func TestNewState_MkdirError(t *testing.T) {
	// given
	// By creating a file named 'data', we ensure that os.MkdirAll fails.
	f, err := os.Create("data")
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, f.Close())
		assert.NoError(t, os.Remove("data"))
	})

	// when
	state, err := NewState()

	// then
	assert.Nil(t, state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create data directory")
}

func TestNewState_BboltOpenError(t *testing.T) {
	// given
	// Pre-lock the database file to cause a timeout
	assert.NoError(t, os.MkdirAll("./data", 0755))
	db, err := bbolt.Open("./data/aiload.bbolt", 0600, nil)
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, os.RemoveAll("./data"))
	})

	// when
	state, err := NewState()

	// then
	assert.Nil(t, state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open bbolt cache")
}
