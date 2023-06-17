package postgres

import (
	"context"
	"fmt"
	"homework-7/internal/pkg/db"
	"homework-7/internal/tests/config"
	"strings"
	"sync"
	"testing"
)

type TDB struct {
	sync.Mutex
	DB *db.Database
}

func NewFromEnv() *TDB {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	database, err := db.NewDB(context.Background(), psqlConn)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &TDB{DB: database}
}

func (d *TDB) SetUp(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	d.Lock()
	d.Truncate(ctx)
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	d.Truncate(context.Background())
}

func (d *TDB) Truncate(ctx context.Context) {
	var tables []string

	err := d.DB.Select(context.Background(), &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(tables) == 0 {
		panic("no migrations")
	}
	q := fmt.Sprintf("Truncate table %s", strings.Join(tables, ","))
	if _, err := d.DB.Exec(context.Background(), q); err != nil {
		panic(err)
	}

	//Тут я обнуляю сиквенсы чтобы тесты нормально работали и можно было прогонять их сколько хочешь
	//(подскажи норм вариант как сделать по другому)
	if _, err := d.DB.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART"); err != nil {
		panic(err)
	}
	if _, err := d.DB.Exec(context.Background(), "ALTER SEQUENCE cars_id_seq RESTART"); err != nil {
		panic(err)
	}
}
