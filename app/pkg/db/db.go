package db

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "time"

    _ "github.com/lib/pq"
)

var DB *sql.DB

// Connect opens a connection to Postgres using DATABASE_URL or individual env vars.
func Connect() error {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        user := os.Getenv("POSTGRES_USER")
        pass := os.Getenv("POSTGRES_PASSWORD")
        name := os.Getenv("POSTGRES_DB")
        host := os.Getenv("POSTGRES_HOST")
        if host == "" {
            host = "db"
        }
        port := os.Getenv("POSTGRES_PORT")
        if port == "" {
            port = "5432"
        }
        dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
    }

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return err
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        _ = db.Close()
        return err
    }

    DB = db
    return nil
}

// Get returns the connected *sql.DB (may be nil if not connected)
func Get() *sql.DB {
    return DB
}
