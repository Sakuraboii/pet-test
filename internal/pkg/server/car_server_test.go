package server

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-7/internal/pkg/repository"
	"homework-7/internal/tests/fixtures"
	"net/http"
	"net/url"
	"testing"
)

var (
	//Машина которого мы хотим создать
	car        = fixtures.Car().Model("Duster").UserId(1).Pointer()
	jsonCar, _ = json.Marshal(car)
)

func Test_unit_getCar(t *testing.T) {
	var (
		id = 1
	)
	t.Run("success,db", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodGet, "car?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		s.carsRepo.EXPECT().GetById(gomock.Any(), int64(id)).Return(&repository.Car{ID: 1, Model: "asd"}, nil)

		// act
		_, status := s.carsServer.getCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		t.Parallel()
		tt := []struct {
			name    string
			request *url.URL
			isOk    bool
		}{
			{
				"without id",
				&url.URL{RawQuery: "car?id"},
				false,
			},
			{
				"wrong id",
				&url.URL{RawQuery: "car?id=asdasd"},
				false,
			},
			{
				"empty",
				&url.URL{RawQuery: ""},
				false,
			},
			{
				"ok",
				&url.URL{RawQuery: "car?id=1"},
				true,
			},
		}
		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				id, err := getCarID(tc.request)
				if !tc.isOk {
					assert.EqualError(t, err, "can't get id")
				} else {
					assert.Equal(t, uint64(0), id)
				}
			})
		}

	})
}

func Test_unit_createCar(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodPost, "car", bytes.NewReader(jsonCar))
		require.NoError(t, err)

		s.carsRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(int64(1), nil)

		// act
		_, status := s.carsServer.createCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodPost, "car", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := s.carsServer.createCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func Test_unit_deleteCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodDelete, "car?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		s.carsRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

		// act
		status := s.carsServer.deleteCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodDelete, "car?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		s.carsRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(assert.AnError)

		// act
		status := s.carsServer.deleteCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func Test_unit_updateCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodPut, "car", bytes.NewReader(jsonCar))
		require.NoError(t, err)

		s.carsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(true, nil)

		// act
		status := s.carsServer.updateCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpCars(t)
		defer s.tearDownCars()

		req, err := http.NewRequest(http.MethodPut, "car", bytes.NewReader(jsonCar))
		require.NoError(t, err)

		s.carsRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(false, assert.AnError)

		// act
		status := s.carsServer.updateCar(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}
