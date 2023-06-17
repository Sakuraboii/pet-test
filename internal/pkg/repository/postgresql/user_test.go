package postgresql

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-7/internal/tests/fixtures"
	"testing"
)

func TestUsersRepo_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)

		//act
		result, err := userRepo.Add(ctx, "Kostya")

		//assert
		assert.NoError(t, err)
		assert.Equal(t, int(result), 1)
	})
}

func TestUsersRepo_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)
		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		//act
		result, err := userRepo.GetById(ctx, 1)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, int(result.ID), 1)
		assert.Equal(t, result.Name, "Kostya")
	})
	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)
		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		//act
		_, err = userRepo.GetById(ctx, 2)

		//assert
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "object not found")
	})
}

func TestUsersRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)
		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		//act
		err = userRepo.Delete(ctx, 1)

		assert.NoError(t, err)
	})
}

func TestUsersRepo_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)
		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		//act
		users, err := userRepo.List(ctx)

		assert.NoError(t, err)
		assert.Equal(t, len(users), 1)
		assert.Equal(t, users[0].Name, "Kostya")
		assert.Equal(t, int(users[0].ID), 1)
	})
}

func TestUsersRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)
		id, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		user := fixtures.User().Name("Dima").Id(id).Pointer()

		//act
		updated, err := userRepo.Update(ctx, user)
		updatedUser, err := userRepo.GetById(ctx, 1)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, updated, true)
		assert.Equal(t, updatedUser.Name, "Dima")
	})
	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		userRepo := NewUsers(Db.DB)
		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		user := fixtures.User().Name("Dima").Id(2).Pointer()

		//act
		updated, err := userRepo.Update(ctx, user)

		//assert
		assert.Error(t, err)
		assert.Equal(t, updated, false)
		assert.Equal(t, err.Error(), "object not found")
	})
}
