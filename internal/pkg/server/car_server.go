package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework-7/internal/pkg/repository"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type serverCarData struct {
	ID     uint64 `json:"ID"`
	Model  string `json:"model"`
	UserID int64  `json:"userID"`
}

// NewCarServerObject Эта функция нужна для получения экзмепляра сервера
func NewCarServerObject(carRepo repository.CarsRepo) *CarServer {
	return &CarServer{
		carRepo: carRepo,
	}
}

type CarServer struct {
	carRepo repository.CarsRepo
}

func (s *CarServer) getCar(cxt context.Context, req *http.Request) ([]byte, int) {
	id, err := getCarID(req.URL)
	if err != nil {
		fmt.Errorf("can't parse id: %s", err)
		return nil, http.StatusBadRequest
	}
	var car *repository.Car
	car, err = s.carRepo.GetById(cxt, int64(id))
	if err != nil {
		fmt.Errorf("can't parse id: %s", err)
		return nil, http.StatusInternalServerError

	}

	su := &serverCarData{}
	su.ID = uint64(car.ID)
	su.Model = car.Model
	su.UserID = car.UserId

	data, err := json.Marshal(su)
	if err != nil {
		fmt.Errorf("can't marshal car with id: %d. Error: %s", id, err)
		return nil, http.StatusInternalServerError
	}

	return data, http.StatusOK
}

func (s *CarServer) createCar(cxt context.Context, req *http.Request) (uint, int) {
	car, err := getCarData(req.Body)
	if err != nil {
		fmt.Println(err)
		return 0, http.StatusBadRequest
	}

	id, err := s.carRepo.Add(cxt, &repository.Car{
		Model:  car.Model,
		UserId: car.UserID,
	})
	if err != nil {
		fmt.Println(err)
		return 0, http.StatusInternalServerError
	}

	return uint(id), http.StatusOK
}

func (s *CarServer) updateCar(ctx context.Context, req *http.Request) int {
	car, err := getCarData(req.Body)
	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest
	}

	u := &repository.Car{ID: int64(car.ID), Model: car.Model}

	updated, err := s.carRepo.Update(ctx, u)

	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest
	}

	if !updated {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (s *CarServer) deleteCar(ctx context.Context, req *http.Request) int {
	id, err := getCarID(req.URL)
	if err != nil {
		fmt.Println(err)
		return http.StatusBadRequest
	}

	e := s.carRepo.Delete(ctx, int64(id))
	if e != nil {
		fmt.Println(err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// CreateCarServer Эта функция для создания хендла машин
func CreateCarServer(ctx context.Context, cr repository.CarsRepo) *http.ServeMux {
	serv := CarServer{
		carRepo: cr,
	}
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/cars", func(res http.ResponseWriter, req *http.Request) {
		if req == nil {
			return
		}

		switch req.Method {
		case http.MethodGet:
			data, status := serv.getCar(ctx, req)
			res.WriteHeader(status)
			res.Write(data)
		case http.MethodPost:
			_, status := serv.createCar(ctx, req)
			res.WriteHeader(status)

		default:
			fmt.Printf("unsupported method: [%s]", req.Method)
			res.WriteHeader(http.StatusNotImplemented)
		}
	})

	return serveMux
}

func getCarData(reader io.ReadCloser) (serverCarData, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return serverCarData{}, err
	}

	data := serverCarData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getCarID(reqUrl *url.URL) (uint64, error) {
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
