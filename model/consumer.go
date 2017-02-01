package model

import(
  "log"
  "fmt"
  "net/url"
  "net/http"

  "github.com/mrjones/oauth"
)

var tokens map[string]*oauth.RequestToken
var c *oauth.Consumer
var client *http.Client = nil

/* Crea un objeto Consumer de la libreria oauth, y un mapa de string y *RequestToken (creo que lo puedo eliminar)
 * para que se utilizen en el controlador.
 */
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

// Una vez creado un cliente nuevo se lo setea y se retorna el mismo para usarlo luego (creo que tambien esta al pedo)
func SetClient(cli *http.Client)(*http.Client) {
  client = cli
  return client
}

/* Con el token y secret que vienen en la url puedo crear un ouath.AccessToken
 * que luego de obtener un consumer con GetConsumer() puedo crear un http.Client
 * que lo retorno. En caso de que falten parametros o sean incorrectos, el metodo retorna nil
 */
func GetClient(v url.Values) *http.Client {
  c,_ := GetConsumer()

  tk := v.Get("token")
  sc := v.Get("secret")
  if tk == "" || sc == "" {
    fmt.Println("\n\n\n FALTAN PARAMETROS \n\n\n")
    return nil
  }
  var accessToken oauth.AccessToken
  accessToken.Token = tk
  accessToken.Secret = sc

  client, err := c.MakeHttpClient(&accessToken)
  if err != nil {
    log.Fatal(err)
  }

  return client
}
