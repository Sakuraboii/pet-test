package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"homework-7/internal/pkg/repository"
)

type serverUserData struct {
	ID   uint64
	Name string
}

type UserServer struct {
	userRepo repository.UsersRepo
}

// NewUserServerObject Эта функция нужна для получения экзмепляра сервера
func NewUserServerObject(userRepo repository.UsersRepo) *UserServer {
	return &UserServer{
		userRepo: userRepo,
	}
}

func (s *UserServer) getUser(cxt context.Context, req *http.Request) ([]byte, int) {
	id, err := getUserID(req.URL)
	if err != nil {
		fmt.Errorf("can't parse id: %s", err)
		return nil, http.StatusBadRequest
	}
	var user *repository.User

	user, err = s.userRepo.GetById(cxt, int64(id))
	if err != nil {
		fmt.Errorf("can't parse id: %s", err)
		return nil, http.StatusInternalServerError

	}

	su := &serverUserData{}
	su.ID = uint64(user.ID)
	su.Name = user.Name

	data, err := json.Marshal(su)
	if err != nil {
		fmt.Errorf("can't marshal user with id: %d. Error: %s", id, err)
		return nil, http.StatusInternalServerError
	}

	return data, http.StatusOK
}

func (s *UserServer) createUser(cxt context.Context, req *http.Request) (uint, int) {
	user, err := getUserData(req.Body)
	if err != nil {
		fmt.Println(err)
		return 0, http.StatusBadRequest
	}
	id, err := s.userRepo.Add(cxt, user.Name)
	if err != nil {
		fmt.Println(err)
		return 0, http.StatusInternalServerError
	}

	return uint(id), http.StatusOK
}

func (s *UserServer) updateUser(ctx context.Context, req *http.Request) int {
	user, err := getUserData(req.Body)
	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest
	}

	u := &repository.User{ID: int64(user.ID), Name: user.Name}

	updated, err := s.userRepo.Update(ctx, u)

	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest
	}

	if !updated {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (s *UserServer) deleteUser(ctx context.Context, req *http.Request) int {
	id, err := getUserID(req.URL)
	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest
	}

	e := s.userRepo.Delete(ctx, int64(id))
	if e != nil {
		fmt.Println(err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// CreateUserServer Эта функция для создания хендла пользователей
func CreateUserServer(ctx context.Context, ur repository.UsersRepo) *http.ServeMux {
	serv := UserServer{
		userRepo: ur,
	}
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/user", func(res http.ResponseWriter, req *http.Request) {
		if req == nil {
			return
		}

		switch req.Method {
		case http.MethodGet:
			data, status := serv.getUser(ctx, req)
			res.WriteHeader(status)
			res.Write(data)
		case http.MethodPost:
			_, status := serv.createUser(ctx, req)
			res.WriteHeader(status)
		case http.MethodDelete:
			status := serv.deleteUser(ctx, req)
			res.WriteHeader(status)
		case http.MethodPut:
			status := serv.updateUser(ctx, req)
			res.WriteHeader(status)

		default:
			fmt.Printf("unsupported method: [%s]", req.Method)
			res.WriteHeader(http.StatusNotImplemented)
		}
	})

	return serveMux
}

func getUserData(reader io.ReadCloser) (serverUserData, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return serverUserData{}, err
	}

	data := serverUserData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getUserID(reqUrl *url.URL) (uint64, error) {
	idStr := reqUrl.Query().Get("id")
	if len(idStr) == 0 {
		return 0, errors.New("can't get id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("can't parse id: %s", err)
	}

	return uint64(id), nil
}
