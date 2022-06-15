package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func initializePostgres(database Database) (*sql.DB, error) {
	d, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			database.Host,
			database.Port,
			database.Username,
			database.Password,
			database.DBName,
		),
	)
	if err != nil {
		log.Panicln(err)
	}

	//d.SetMaxOpenConns(database.MaxOpenConnections)
	//d.SetMaxIdleConns(database.MaxIdleConnections)

	if err := d.Ping(); err != nil {
		return nil, err
	}

	return d, nil
}
