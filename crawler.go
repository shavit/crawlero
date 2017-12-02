package crawlero

import (
  "errors"
  "fmt"
  "io/ioutil"
  "log"
  "net/url"
  "net/http"
  "os"
  "strings"
  "time"

  "golang.org/x/net/proxy"
  "github.com/streadway/amqp"
)

type Crawler interface {

  // SetProxy set a proxy URL
  SetProxy(u string) (err error)

  // Save opens and saves a webpage
  Save(u string, retry int) (err error)

  // Listen listens to incoming messages
  Listen(done chan error)
}

type crawler struct {
  conn *Connection
  httpClient *http.Client
  page string
}

func NewCrawler(conn *Connection) *crawler {
  return &crawler{
    conn: conn,
    httpClient: new(http.Client),
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

  cw.httpClient.Transport = &http.Transport{Dial: dialer.Dial}

  return err
}

// Save opens and saves a webpage
func (cw *crawler) Save(u string, retry int) (err error) {
  if retry < 0 {
    return errors.New("Too many attempts")
  }

  var resp *http.Response
  var respBody []byte
  var body string
  resp, err = cw.httpClient.Get(u)
  if err != nil && retry == 0 {
    return err
  } else if err != nil {
    time.Sleep(time.Duration(600 / retry * retry) * time.Second)
    cw.Save(u, retry-1)
    return
  }
  defer resp.Body.Close()

  respBody, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  body = string(respBody)
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

func (cw *crawler) Listen(done chan error) {
  var err error
  var user string = os.Getenv("RABBITMQ_DEFAULT_USER")
  var pass string = os.Getenv("RABBITMQ_DEFAULT_PASS")
  var vhost string = os.Getenv("RABBITMQ_DEFAULT_VHOST")
  var host string = os.Getenv("RABBITMQ_DEFAULT_HOST")
  var connection *amqp.Connection
  var channel *amqp.Channel
  var queue amqp.Queue
  var messages <-chan amqp.Delivery = make(<-chan amqp.Delivery)

  connection, err = amqp.Dial("amqp://"+user+":"+pass+"@"+host+"/"+vhost)
  if err != nil {
    done <- err
    return
  }
  defer connection.Close()

  channel, err = connection.Channel()
  if err != nil {
    done <- err
    return
  }
  defer channel.Close()

  err = channel.ExchangeDeclare(
    "crawlers", // Name
    "direct", // Type
    true, // Durable
    false, // Auto-deleted
    false, // Internal
    false, // No-wait
    nil, // Arguments
  )
  if err != nil {
    done <- err
    return
  }

  // Declare the queue on both the consumer and publisher, because it might
  //  start before the publisher started
  queue, err = channel.QueueDeclare(
    "save_page", // Name
    true, // Durable
    true, // Delete when unused
    false, // Exclusive
    false, // No-wait
    nil, // Arguments
  )
  if err != nil {
    done <- err
    return
  }

  err = channel.QueueBind(
    queue.Name, // Queue
    "route_crawlers", // Routing key
    "crawlers", // Exchange
    false,
    nil,
  )
  if err != nil {
    done <- err
    return
  }

  messages, err = channel.Consume(
    queue.Name, // Queue
    "", // consumer
    false, // Auto acknoledge
    false, // Exclusive
    false, // No-local
    false, // No-wait
    nil, // Args
  )
  if err != nil {
    done <- err
    return
  }

  for m := range messages {
    log.Println("New message", string(m.Body))
    err = cw.Save(string(m.Body), 10)
    if err == nil {
      m.Ack(true)
    }
  }

  done <- nil
}
