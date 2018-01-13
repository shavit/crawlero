package main

import (
  "fmt"
  "os"
  "strconv"
  "sync"

  "github.com/shavit/crawlero"
)

func main(){
  var input string
  var wg sync.WaitGroup
  var done chan error

  println(`
    How many workers to start?
  `)
  fmt.Scanln(&input)

  n, _ := strconv.Atoi(input)
  if n > 0 && n < 40 {
    wg.Add(n)
    dbConn := crawlero.NewConnection()

    for i := 0; i < n; i++ {
      go func(i int, done chan error) {
        defer wg.Done()
        done = make(chan error)
        cw := crawlero.NewCrawler(dbConn)
        err := cw.SetProxy("socks5://127.0.0.1:9050")
        if err != nil {
          panic(err)
        }
        go cw.Listen(done)
        println("Starting worker", i)
        <-done
      }(i, done)
    }
    println("Waiting for", n, "workers")
    wg.Wait()
  } else {
    println("Invalid input")
    os.Exit(1)
  }
}
