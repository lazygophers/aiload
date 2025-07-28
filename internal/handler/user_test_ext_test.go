package handler_test

import (
	"aiload/internal/handler"
	"aiload/internal/model"
	"aiload/internal/repository"
	"aiload/internal/server"
	"aiload/internal/service"
	"aiload/internal/state"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// setupApp initializes a new Fiber app with test-specific dependencies.
func setupApp(t *testing.T) (*fiber.App, func()) {
	t.Helper()

	db, dbCleanup := state.NewTestDB(t)
	cache, cacheCleanup := state.NewTestCache(t)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cache)
	userHandler := handler.NewUserHandler(userService)

	// We need a way to pass the handler to the routes
	app := fiber.New()
	server.SetupRoutes(app, userHandler)

	cleanup := func() {
		dbCleanup()
		cacheCleanup()
	}

	return app, cleanup
}

func TestUserHandler(t *testing.T) {
	app, cleanup := setupApp(t)
	defer cleanup()

	var createdUser model.User

	t.Run("CreateUser", func(t *testing.T) {
		user := &model.User{Name: "Test User", Email: "test@example.com"}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/api/users/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1) // -1 for no timeout
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&createdUser)
		assert.NoError(t, err)
		assert.NotZero(t, createdUser.ID)
		assert.Equal(t, user.Name, createdUser.Name)
	})

	t.Run("GetUser", func(t *testing.T) {
		if createdUser.ID == 0 {
			t.Skip("Skipping GetUser test as CreateUser failed")
		}
		req := httptest.NewRequest(http.MethodGet, "/api/users/"+strconv.Itoa(int(createdUser.ID)), nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var fetchedUser model.User
		err = json.NewDecoder(resp.Body).Decode(&fetchedUser)
		assert.NoError(t, err)
		assert.Equal(t, createdUser.ID, fetchedUser.ID)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		if createdUser.ID == 0 {
			t.Skip("Skipping UpdateUser test as CreateUser failed")
		}
		updatedUser := &model.User{Name: "Updated User", Email: "updated@example.com"}
		body, _ := json.Marshal(updatedUser)
		req := httptest.NewRequest(http.MethodPut, "/api/users/"+strconv.Itoa(int(createdUser.ID)), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var returnedUser model.User
		err = json.NewDecoder(resp.Body).Decode(&returnedUser)
		assert.NoError(t, err)
		assert.Equal(t, "Updated User", returnedUser.Name)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		if createdUser.ID == 0 {
			t.Skip("Skipping DeleteUser test as CreateUser failed")
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/users/"+strconv.Itoa(int(createdUser.ID)), nil)

		resp, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify user is deleted
		req = httptest.NewRequest(http.MethodGet, "/api/users/"+strconv.Itoa(int(createdUser.ID)), nil)
		resp, err = app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("ErrorCases", func(t *testing.T) {
		t.Run("CreateUser_InvalidJSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/users/", bytes.NewReader([]byte(`{"invalid`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("GetUser_InvalidID", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/users/abc", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("GetUser_NotFound", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/users/99999", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

		t.Run("UpdateUser_InvalidID", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/users/abc", bytes.NewReader([]byte(`{}`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("UpdateUser_InvalidJSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewReader([]byte(`{"invalid`)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("DeleteUser_InvalidID", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/users/abc", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
