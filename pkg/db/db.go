package db

import (
	"fmt"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Init(c config.DbConfig) (*sqlx.DB, error) {

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Pass, c.Name, c.SllMode,
	)

	db, err := sqlx.Connect("postgres", connection)

	if err != nil {
		return nil, err
	}

	return db, err
}
