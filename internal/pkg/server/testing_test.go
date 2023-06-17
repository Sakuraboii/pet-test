package server

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_database "homework-7/internal/pkg/db/mocks"
	"homework-7/internal/pkg/repository/mocks"
	"homework-7/internal/tests/postgres"
	"testing"
)

var (
	ctx = context.Background()
	Db  *postgres.TDB
)

type usersRepoFixture struct {
	ctrl       *gomock.Controller
	usersRepo  *mocks.MockUsersRepo
	mockDb     *mock_database.MockDBops
	userServer *UserServer
}

func setUpUsers(t *testing.T) usersRepoFixture {
	ctrl := gomock.NewController(t)

	mockDb := mock_database.NewMockDBops(ctrl)
	repo := mocks.NewMockUsersRepo(ctrl)
	uServer := NewUserServerObject(repo)
	return usersRepoFixture{usersRepo: repo, ctrl: ctrl, mockDb: mockDb, userServer: uServer}

}

func (u *usersRepoFixture) tearDownUsers() {
	u.ctrl.Finish()
}

type carsRepoFixture struct {
	ctrl       *gomock.Controller
	carsRepo   *mocks.MockCarsRepo
	mockDb     *mock_database.MockDBops
	carsServer *CarServer
}

func setUpCars(t *testing.T) carsRepoFixture {
	ctrl := gomock.NewController(t)

	mockDb := mock_database.NewMockDBops(ctrl)
	repo := mocks.NewMockCarsRepo(ctrl)
	cServer := NewCarServerObject(repo)
	return carsRepoFixture{carsRepo: repo, ctrl: ctrl, mockDb: mockDb, carsServer: cServer}

}

func (u *carsRepoFixture) tearDownCars() {
	u.ctrl.Finish()
}

func init() {
	Db = postgres.NewFromEnv()
}
