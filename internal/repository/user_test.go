package repository

import (
	"aiload/internal/model"
	"aiload/internal/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	db, cleanup := state.NewTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create
	user := &model.User{Name: "Test User", Email: "test@example.com"}
	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Get
	foundUser, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", foundUser.Name)

	// Update
	foundUser.Name = "Updated User"
	err = repo.UpdateUser(foundUser)
	assert.NoError(t, err)
	updatedUser, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User", updatedUser.Name)

	// Delete
	err = repo.DeleteUser(user.ID)
	assert.NoError(t, err)
	_, err = repo.GetUserByID(user.ID)
	assert.Error(t, err)
}
