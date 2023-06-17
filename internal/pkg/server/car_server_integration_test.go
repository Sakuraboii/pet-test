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

func Test_integration_getCar(t *testing.T) {
	t.Run("success,db", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		_, err = carRepo.Add(ctx, car)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "car?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := carServer.getCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("fail", func(t *testing.T) {
		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		_, err := carRepo.Add(ctx, car)

		req, err := http.NewRequest(http.MethodGet, "car?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := carServer.getCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func Test_integration_createCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "car", bytes.NewReader(jsonCar))
		require.NoError(t, err)

		// act
		id, status := carServer.createCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, int(id), 1)
	})
	t.Run("fail", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		req, err := http.NewRequest(http.MethodPost, "car", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := carServer.createCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func Test_integration_deleteCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		_, err = carRepo.Add(ctx, car)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, "car?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		status := carServer.deleteCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
}

func Test_integration_updateCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := postgresql.NewUsers(Db.DB)
		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		_, err := userRepo.Add(ctx, "Kostya")
		require.NoError(t, err)

		_, err = carRepo.Add(ctx, car)
		require.NoError(t, err)

		updatedCar := fixtures.Car().Model("Logan").Id(1).Pointer()
		jsonUpdatedCar, _ := json.Marshal(updatedCar)

		req, err := http.NewRequest(http.MethodPut, "car", bytes.NewReader(jsonUpdatedCar))
		require.NoError(t, err)

		// act
		status := carServer.updateCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {

		// arrange
		Db.SetUp(t)
		defer Db.TearDown()

		carRepo := postgresql.NewCars(Db.DB)
		carServer := NewCarServerObject(carRepo)

		updatedCar := fixtures.Car().Model("Logan").Id(1).Pointer()
		jsonUpdatedCar, _ := json.Marshal(updatedCar)

		req, err := http.NewRequest(http.MethodPut, "car", bytes.NewReader(jsonUpdatedCar))
		require.NoError(t, err)

		// act
		status := carServer.updateCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}
