package crawlero

import (
  "testing"
)

func TestConnectionCanConnect(t *testing.T){
  var conn *Connection = NewConnection()

  if conn == nil {
    t.Error("Got nil while expecting a connection")
  }

  if conn.DB == nil {
    t.Error("Got nil while expecting DB connection")
  }
}
