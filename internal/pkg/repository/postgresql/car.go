package postgresql

import (
	"context"
	"database/sql"
	"homework-7/internal/pkg/db"
	"homework-7/internal/pkg/repository"
)

type CarsRepo struct {
	db db.DBops
}

func NewCars(db db.DBops) *CarsRepo {
	return &CarsRepo{db: db}
}

func (r *CarsRepo) Add(ctx context.Context, car *repository.Car) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO cars(model,user_id) VALUES ($1,$2) RETURNING id`,
		car.Model, car.UserId).Scan(&id)
	return id, err
}

func (r *CarsRepo) GetById(ctx context.Context, id int64) (*repository.Car, error) {
	var car repository.Car
	err := r.db.Get(ctx, &car, "SELECT id,model,user_id,created_at,updated_at FROM cars WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNotFound
	}
	return &car, err
}

func (r *CarsRepo) List(ctx context.Context) ([]*repository.Car, error) {
	cars := make([]*repository.Car, 0)
	err := r.db.Select(ctx, &cars, "SELECT id, model, user_id, created_at, updated_at FROM cars")
	if err != nil {
		return nil, repository.ErrObjectNotFound
	}
	return cars, err
}

func (r *CarsRepo) Update(ctx context.Context, car *repository.Car) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE cars SET model = $1 WHERE id = $2", car.Model, car.ID)
	if err != nil || result.RowsAffected() == 0 {
		return false, repository.ErrObjectNotFound
	}
	return result.RowsAffected() > 0, nil
}

func (r *CarsRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
