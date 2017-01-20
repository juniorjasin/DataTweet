package model

import(
  "github.com/mrjones/oauth"
)

var tokens map[string]*oauth.RequestToken
var c *oauth.Consumer

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
