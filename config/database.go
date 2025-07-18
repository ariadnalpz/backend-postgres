package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Conn *pgxpool.Pool

func InitDatabase() error {
	connStr := os.Getenv("SUPABASE_CONNECTION_STRING")
	var err error
	Conn, err = pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		return err
	}
	return nil
}
