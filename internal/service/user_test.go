package service

import (
	"aiload/internal/model"
	"aiload/internal/repository"
	"aiload/internal/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	db, dbCleanup := state.NewTestDB(t)
	defer dbCleanup()

	cache, cacheCleanup := state.NewTestCache(t)
	defer cacheCleanup()

	repo := repository.NewUserRepository(db)
	service := NewUserService(repo, cache)

	// Create
	user := &model.User{Name: "Test User", Email: "test@example.com"}
	err := service.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Get from DB (cache miss)
	foundUser, err := service.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", foundUser.Name)

	// Get from Cache (cache hit)
	cachedUser, err := service.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, foundUser.ID, cachedUser.ID)
	assert.Equal(t, foundUser.Name, cachedUser.Name)
	assert.Equal(t, foundUser.Email, cachedUser.Email)

	// Update (and invalidate cache)
	foundUser.Name = "Updated User"
	err = service.UpdateUser(foundUser)
	assert.NoError(t, err)

	// Get again (should be from DB, since cache was invalidated)
	updatedUser, err := service.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User", updatedUser.Name)

	// Delete
	err = service.DeleteUser(user.ID)
	assert.NoError(t, err)
	_, err = service.GetUser(user.ID)
	assert.Error(t, err)
}
