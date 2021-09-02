package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	// required dependencie for migrate
	_ "github.com/golang-migrate/migrate/source/file"
	// postgres driver
	_ "github.com/lib/pq"
)

var instance *sql.DB

// GetDB returns the database single instance
func GetDB() *sql.DB {
	if instance == nil {
		log.Print("creating instance of repository")
		instance = initDb()
		migrateDB()
	}
	return instance
}

func initDb() *sql.DB {
	log.Println("connecting to db")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PWD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not create connection to DB, %v", err)
	}

	connPoolSize, err := strconv.Atoi(os.Getenv("CONN_POOL_SIZE"))
	if err != nil {
		panic(fmt.Errorf("could not crate the connection pool, %v", err))
	}

	db.SetMaxIdleConns(connPoolSize)
	db.SetMaxOpenConns(connPoolSize)
	db.SetConnMaxLifetime(120 * time.Minute)
	log.Println("connected to db")
	return db
}

func migrateDB() {
	driver, err := postgres.WithInstance(instance, &postgres.Config{})
	if err != nil {
		log.Printf("could not migrate db, %v", err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/database/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Printf("could not migrate db, %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			log.Printf("could not migrate db, %v", err)
		}
	}
}

// Close closes the database connection
func Close() {
	instance.Close()
}
