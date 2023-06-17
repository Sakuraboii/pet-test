package server

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-7/internal/pkg/repository/postgresql"
	"homework-7/internal/tests/fixtures"
	"net/http"
	"testing"
)

func Test_integration_getUser(t *testing.T) {
	t.Run("success,db", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "user?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := userServer.getUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		req, err := http.NewRequest(http.MethodGet, "user?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := userServer.getUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func Test_integration_createUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		req, err := http.NewRequest(http.MethodPost, "user", bytes.NewReader(jsonUser))
		require.NoError(t, err)

		// act
		id, status := userServer.createUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, int(id), 1)
	})
	t.Run("fail", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		req, err := http.NewRequest(http.MethodPost, "user", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := userServer.createUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func Test_integration_deleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, "user?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		res := userServer.deleteUser(ctx, req)

		// assert
		assert.Equal(t, res, http.StatusOK)
	})
}

func Test_integration_updateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		updatedUser := fixtures.User().Name("Dima").Id(1).Pointer()
		jsonUpdatedUser, _ := json.Marshal(updatedUser)

		req, err := http.NewRequest(http.MethodPut, "user", bytes.NewReader(jsonUpdatedUser))
		require.NoError(t, err)

		// act
		status := userServer.updateUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		userServer := NewUserServerObject(userRepo)

		updatedUser := fixtures.User().Name("Dima").Id(1).Pointer()
		jsonUpdatedUser, _ := json.Marshal(updatedUser)

		req, err := http.NewRequest(http.MethodPut, "user", bytes.NewReader(jsonUpdatedUser))
		require.NoError(t, err)

		// act
		status := userServer.updateUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}
