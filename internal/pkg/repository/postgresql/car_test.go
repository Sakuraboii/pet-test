package postgresql

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-7/internal/tests/fixtures"
	"testing"
)

func TestCarsRepo_Add(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)
		userRepo := NewUsers(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()
		_, err := userRepo.Add(ctx, "Kostya")

		//act
		result, err := carsRepo.Add(ctx, car)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, int(result), 1)
	})
	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()

		//act
		result, err := carsRepo.Add(ctx, car)

		//assert
		assert.Error(t, err)
		assert.Equal(t, int(result), 0)
	})
}

func TestCarsRepo_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)
		userRepo := NewUsers(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()

		_, err := userRepo.Add(ctx, "Kostya")
		_, er := carsRepo.Add(ctx, car)
		require.NoError(t, er)

		//act
		result, err := carsRepo.GetById(ctx, 1)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, int(result.ID), 1)
		assert.Equal(t, result.Model, "Duster")
	})
	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)

		//act
		result, err := carsRepo.GetById(ctx, 1)

		//assert
		assert.Error(t, err)
		assert.Equal(t, int(result.ID), 0)
		assert.Equal(t, result.Model, "")
	})
}

func TestCarsRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)
		userRepo := NewUsers(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()

		_, err := userRepo.Add(ctx, "Kostya")
		_, err = carsRepo.Add(ctx, car)
		require.NoError(t, err)

		//act
		err = carsRepo.Delete(ctx, 1)

		//assert
		assert.NoError(t, err)
	})
}

func TestCarsrRepo_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)
		userRepo := NewUsers(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()

		_, err := userRepo.Add(ctx, "Kostya")
		_, err = carsRepo.Add(ctx, car)
		require.NoError(t, err)

		//act
		cars, err := carsRepo.List(ctx)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, len(cars), 1)
		assert.Equal(t, cars[0].Model, "Duster")
		assert.Equal(t, int(cars[0].UserId), 1)
		assert.Equal(t, int(cars[0].ID), 1)
	})
}

func TestCarsRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)
		userRepo := NewUsers(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()

		_, err := userRepo.Add(ctx, "Kostya")
		id, err := carsRepo.Add(ctx, car)

		updateCar := fixtures.Car().Model("Logan").Id(id).Pointer()

		require.NoError(t, err)

		//act
		updated, err := carsRepo.Update(ctx, updateCar)
		updatedCar, err := carsRepo.GetById(ctx, 1)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, updated, true)
		assert.Equal(t, updatedCar.Model, "Logan")
	})
	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		carsRepo := NewCars(Db.DB)
		userRepo := NewUsers(Db.DB)

		car := fixtures.Car().Model("Duster").UserId(1).Pointer()

		_, err := userRepo.Add(ctx, "Kostya")
		_, err = carsRepo.Add(ctx, car)

		updateCar := fixtures.Car().Model("Logan").Id(2).Pointer()

		require.NoError(t, err)

		//act
		updated, err := carsRepo.Update(ctx, updateCar)

		//assert
		assert.Error(t, err)
		assert.Equal(t, updated, false)
		assert.Equal(t, err.Error(), "object not found")
	})
}
