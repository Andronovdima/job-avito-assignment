package apiserver

import (
	"database/sql"
	"fmt"
	"github.com/Andronovdima/job-avito-assignment/store"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func Start() error {
	config := NewConfig()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":9000"
	} else {
		port = ":" + port
	}
	config.BindAddr = port

	url := os.Getenv("DATABASE_URL")
	if len(url) != 0 {
		config.DatabaseURL = url
	}

	zapLogger, _ := zap.NewProduction()
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			log.Println("zap logger error", err)
		}
	}()
	sugaredLogger := zapLogger.Sugar()

	srv, err := NewServer(config, sugaredLogger)
	if err != nil {
		return err
	}

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		sugaredLogger.Error(err)
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			sugaredLogger.Error(err)
		}
	}()

	srv.ConfigureServer(db)

	fmt.Println("Start server on port " + config.BindAddr)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	if err := store.CreateTables(db); err != nil {
		return nil, err
	}
	return db, nil
}
