package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	dsn := "postgres://postgres:258456@localhost:5432/calculator"

	var err error
	DB, err = pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	err = DB.Ping(context.Background())

	if err != nil {
		log.Fatalf("Не удалось выполнить ping БД: %v", err)
	}

	fmt.Println("Подключение к PostgreSQL успешно установлено")

}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
