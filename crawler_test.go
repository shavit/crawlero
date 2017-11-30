package crawlero

import (
  "testing"
  "os"
)

func TestCreateCrawler(t *testing.T){
  var cw Crawler = NewCrawler(NewConnection())

  if cw == nil {
    t.Error("Found nil whlie creating a new crawler")
  }
}

func TestSetProxy(t *testing.T){
  var err error
  var cw Crawler = NewCrawler(NewConnection())

  cw.SetProxy("socks5://127.0.0.1:9050")
  if err != nil {
    t.Error("Error setting a proxy:", err)
  }
}

func TestSave(t *testing.T){
  var err error
  var cw Crawler = NewCrawler(NewConnection())

  err = cw.Save("http://localhost:0/", -1)
  if err.Error() != "Too many attempts" {
    t.Error("Should raise an error after too many attempts.", err)
  }

  err = cw.Save("http://localhost:0/", 0)
  if err.Error() != "Get http://localhost:0/: dial tcp 127.0.0.1:0: getsockopt: connection refused" {
    t.Error(err)
  }
}

func TestListen(t *testing.T){
  var done chan error = make(chan error)
  var cw Crawler = NewCrawler(NewConnection())
  os.Setenv("RABBITMQ_DEFAULT_USER", "wrong-user")

  go cw.Listen(done)
  if <-done == nil {
    t.Error("Should not be able to listen with invalid credentials")
  }
}
