package crawlero

import (
  "os"
  "database/sql"
  _ "github.com/lib/pq"
)

type Connection struct {
  DB *sql.DB
}

func NewConnection() (c *Connection) {
  var err error
  var uri string = os.Getenv("DATABASE_URL")

  c = &Connection{
    DB: new(sql.DB),
  }

  c.DB, err = sql.Open("postgres", uri+"?sslmode=disable")
  if err != nil {
    panic(err)
  }

  return c
}
