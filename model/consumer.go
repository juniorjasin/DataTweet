package model

import(
  "net/http"
  "github.com/mrjones/oauth"
)

var tokens map[string]*oauth.RequestToken
var c *oauth.Consumer
var client *http.Client = nil

func GetConsumer() (*oauth.Consumer,map[string]*oauth.RequestToken) {
  tokens = make(map[string]*oauth.RequestToken)

  var ck string = "8ooc5TaVYDvTjNbIR3zTo72zh"
  var cs string = "HtPRUQJulbqKX79MxASz4feoKq9dnByfRCl6AJHhAUYdgRangx"
  var consumerKey *string = &ck
  var consumerSecret *string = &cs

  c = oauth.NewConsumer(
    *consumerKey,
    *consumerSecret,
    oauth.ServiceProvider{
      RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
      AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
      AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
    },
  )
  c.Debug(true)

  return c,tokens
}

func SetClient(cli *http.Client)(*http.Client) {
  client = cli
  return client
}

func GetClient() (*http.Client){
  if client != nil{
    return client
  }else{
    return nil
  }
}

// tipos y metodos para ordenar mapa
type Pair struct {
  Key string
  Value float64
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }
