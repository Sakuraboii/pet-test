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
	//Пользователь которого мы хотим создать
	user        = fixtures.User().Name("Kostya").Pointer()
	jsonUser, _ = json.Marshal(user)
)

func Test_unit_getUser(t *testing.T) {
	var (
		id = 1
	)
	t.Run("success,db", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodGet, "user?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		s.usersRepo.EXPECT().GetById(gomock.Any(), int64(id)).Return(&repository.User{ID: 1, Name: "asd"}, nil)

		// act
		_, status := s.userServer.getUser(ctx, req)

		// assert
		require.Equal(t, http.StatusOK, status)
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
				&url.URL{RawQuery: "user?id"},
				false,
			},
			{
				"wrong id",
				&url.URL{RawQuery: "user?id=asdasd"},
				false,
			},
			{
				"empty",
				&url.URL{RawQuery: ""},
				false,
			},
			{
				"ok",
				&url.URL{RawQuery: "user?id=1"},
				true,
			},
		}
		for _, tc := range tt {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				id, err := getUserID(tc.request)
				if !tc.isOk {
					assert.EqualError(t, err, "can't get id")
				} else {
					assert.Equal(t, uint64(0), id)
				}
			})
		}

	})
}

func Test_unit_createUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodPost, "user", bytes.NewReader(jsonUser))
		require.NoError(t, err)

		s.usersRepo.EXPECT().Add(gomock.Any(), "Kostya").Return(int64(1), nil)

		// act
		_, status := s.userServer.createUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodPost, "user", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		// act
		_, status := s.userServer.createUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func Test_unit_deleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodDelete, "user?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		s.usersRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

		// act
		status := s.userServer.deleteUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodDelete, "user?id=1", bytes.NewReader([]byte{}))
		require.NoError(t, err)

		s.usersRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(assert.AnError)

		// act
		status := s.userServer.deleteUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func Test_unit_updateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodPut, "user", bytes.NewReader(jsonUser))
		require.NoError(t, err)

		s.usersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(true, nil)

		// act
		status := s.userServer.updateUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusOK, status)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUpUsers(t)
		defer s.tearDownUsers()

		req, err := http.NewRequest(http.MethodPut, "user", bytes.NewReader(jsonUser))
		require.NoError(t, err)

		s.usersRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(false, assert.AnError)

		// act
		status := s.userServer.updateUser(ctx, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, status)
	})
}
