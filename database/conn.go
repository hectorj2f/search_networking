package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/hectorj2f/search_networking/resources"

	_ "github.com/lib/pq"
)

var db *DB

type DB struct {
	Conn *sql.DB

	dsnSuffix string

	mtx  sync.RWMutex
	dsn  string
	addr string

	stmts map[string]*sql.Stmt
}

func SetupConnection() (*DB, error){
	var err error
	db, err = connect()
	return db, err
}

func GetDatabase() (*DB) {
	return db
}

func connect() (*DB, error) {
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	database := os.Getenv("PGDATABASE")
	username := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	sslmode := os.Getenv("PGSSLMODE")

	if host == "" {
		host = resources.DB_ADDR
	}
	if port == "" {
		port = resources.DB_PORT
	}
	if database == "" {
		database = resources.DB_NAME
	}
	if username == "" {
		username = resources.DB_USERNAME
	}
	if password == "" {
		password = resources.DB_PASSWORD
	}
	if sslmode == "" {
		sslmode = resources.DB_SSL_MODE
	}

	dsn := "dbname=giantswarm"
	conn, err := sql.Open(resources.DB_SERVICE, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, host, port, database, sslmode))
	postgresDb := &DB{
		Conn:			 conn,
		dsnSuffix: dsn,
		dsn:       fmt.Sprintf("host=leader.%s.discoverd %s", resources.DB_SERVICE, dsn),
		addr:      fmt.Sprintf("leader.%s.discoverd", resources.DB_SERVICE),
		stmts:     make(map[string]*sql.Stmt),
	}

	if err != nil {
		return nil, err
	}

	return postgresDb, nil
}

func (db *DB) Close() {
	db.Conn.Close()
}
