package crawlero

import (
  "errors"
  "fmt"
  "net/url"
  "net/http"
  "strings"

  "golang.org/x/net/proxy"
  "gopkg.in/headzoo/surf.v1"
  "github.com/headzoo/surf/browser"
)

type Crawler interface {
  // TODO: Visit page, save page, set proxy, start multiple workers

  // SetProxy set a proxy URL
  SetProxy(u string) (err error)

  // Save opens and saves a webpage
  Save(u string) (err error)
}

type crawler struct {
  bow *browser.Browser
  conn *Connection
  page string
}

func NewCrawler(conn *Connection) *crawler {
  return &crawler{
    bow: surf.NewBrowser(),
    conn: conn,
  }
}

// SetProxy set a proxy URL
func (cw *crawler) SetProxy(u string) (err error) {
  _url, err := url.Parse(u)
  if err != nil {
    return err
  }

  dialer, err := proxy.FromURL(_url, proxy.Direct)
  if err != nil {
    return err
  }

  cw.bow.SetTransport(&http.Transport{Dial: dialer.Dial})

  return err
}

// Save opens and saves a webpage
func (cw *crawler) Save(u string) (err error) {
  err = cw.bow.Open(u)
  if err != nil {
    return err
  }

  var body string = cw.bow.Body()
  body = strings.Replace(body, "'", "\\'", -1)
  body = strings.Replace(body, "'", "''", -1)
  cw.page = body

  if cw.page == "" {
    return errors.New("Could not save empty webpage")
  }

  var query string = fmt.Sprintf(`INSERT INTO webpages
    (url, body)
    VALUES ('%v', '%v')
    ON CONFLICT (url) DO UPDATE
    SET body = '%v'`, u, cw.page, cw.page)
  _, err = cw.conn.DB.Exec(query)

  return err
}
