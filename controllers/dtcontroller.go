package controllers

import(
  "fmt"
  "log"
  "net/url"
  "net/http"
  "encoding/json"

  "github.com/mrjones/oauth"
  "github.com/juniorjasin/datatweet/model"
  "github.com/juniorjasin/datatweet/services"
)
// variable auxiliar para no perder el valor de tokens
var tokensAux map[string]*oauth.RequestToken
var prueba string

/* Se llama a este metodo cuando llega un request a localhost:8888
 * metodo redirije al usuario para que acepte los permisos
 * siguiendo el estandar oauth1.0 de twitter. Utilizo la libreria "github.com/mrjones/oauth"
 * para realizarlo porque Go solo tiene el package oauth2 para oauth 2.0
 */
func RedirectUserToTwitter(w http.ResponseWriter, r *http.Request) {
  var c *oauth.Consumer
  var tokens map[string]*oauth.RequestToken
  c, tokens = services.GetConsumer()
  tokensAux = tokens
	tokenUrl := fmt.Sprintf("http://%s/maketoken", r.Host)
	token, requestUrl, err := c.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
    err := model.Error{Code: http.StatusServiceUnavailable, Message: "Servicio no disponible. Problemas en la conexion"}
    jsonError, _ := json.Marshal(err)
    fmt.Fprintf(w, string(jsonError))
    return
		log.Fatal(err)
	}
	// Make sure to save the token, we'll need it for AuthorizeToken()
	tokens[token.Token] = token
	http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
}

/* Automaticamente luego de que se aceptaron los permisos, se llama a este metodo,
 * donde se crea el AccessToken, con el que se crea un *http.Client, donde lo guardo
 * para luego mandarle request a la API de twitter
 */
func GetTwitterToken(w http.ResponseWriter, r *http.Request) {
  var c *oauth.Consumer
  var tokens map[string]*oauth.RequestToken
  c, tokens = services.GetConsumer()
  tokens = tokensAux
  values := r.URL.Query()
  accessToken := getAccessToken(values, c, tokens)

  client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
  }

  cli := services.SetClient(client)
  if cli == nil{
    log.Fatal(cli)
  }

  jsonToken, err := json.Marshal(accessToken)
  if err != nil{
    log.Fatal(err)
  }
  // Retorno el json con token, secret y AdditionalData con screen_name del que solicito, user_id y x_auth_expires
  fmt.Fprintf(w, string(jsonToken))
}

// obtengo el access token
func getAccessToken(values url.Values, c *oauth.Consumer, tokens map[string]*oauth.RequestToken) *oauth.AccessToken {
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	accessToken, err := c.AuthorizeToken(tokens[tokenKey], verificationCode)
	if err != nil {
		log.Fatal(err)
	}
  return accessToken
}

// consulta metodo de analisis.go para obtener el porcentaje de favoritos y
// devoluelve la respuesta procesada
func Favorites(w http.ResponseWriter, r *http.Request){
  values := r.URL.Query()
  sn := values.Get("screen_name")
  if sn == "" {
    ResponseError(1,w)
    return
  }

  // obtengo un cliente
  client := services.GetClient(values)
  if client == nil {
    ResponseError(2,w)
    return
  }

  pl,code := services.PercentageOfFavorites(values, client)
  if pl == nil {
    ResponseError(code,w)
    return
  }

  fav := make(services.PairList, 10)
  for i := 0; i < 10; i++ {
      // s64 := strconv.FormatFloat(pl[i].Value, 'f', 2, 64)
      fav[i] = services.Pair{Key: pl[i].Key, Value: pl[i].Value}
    }
  json, _ := json.Marshal(fav)
  fmt.Fprintf(w, string(json))
}

// consulta metodo de dictionary.go para obtener un mapa de palabras y
// porcentaje ordenado
func Dictionary(w http.ResponseWriter, r *http.Request){
  values := r.URL.Query()
  // obtengo un cliente
  client := services.GetClient(values)
  if client == nil {
    ResponseError(2,w)
    return
  }

  sn := values.Get("screen_name")
  if sn == "" {
    ResponseError(1,w)
    return
  }

  pl,code := services.GetDictionary(client, sn)
  if pl == nil{
    ResponseError(code,w)
    return
  }

  dic := make(services.PairList, 0)
  for _,v := range pl {
    dic = append(dic, services.Pair{Key: v.Key, Value: v.Value})
  }

  json, _ := json.Marshal(dic)
  fmt.Fprintf(w, string(json))
}

// metodo que retorna json informando cada tipo de error segun el code que viene
func ResponseError(code float64, w http.ResponseWriter)  {
  switch code {
  case 89:
    err := model.Error{Code: http.StatusBadRequest, Message: "Invalid or expired token."}
    jsonError, _ := json.Marshal(err)
    fmt.Fprintf(w, string(jsonError))
    return

  case 32:
    err := model.Error{Code: http.StatusBadRequest, Message: "Could not authenticate you (problem with secret key)"}
    jsonError, _ := json.Marshal(err)
    fmt.Fprintf(w, string(jsonError))
    return

  case 1:
    err := model.Error{Code: http.StatusBadRequest, Message: "screen_name missing"}
    jsonError, _ := json.Marshal(err)
    fmt.Fprintf(w, string(jsonError))
    return

  case 2:
    err := model.Error{Code: http.StatusBadRequest, Message: "uninitialized client, token or secret missing"}
    jsonError, _ := json.Marshal(err)
    fmt.Fprintf(w, string(jsonError))
    return

  default:
    err := model.Error{Code: http.StatusBadRequest, Message: "Problems with parameters"}
    jsonError, _ := json.Marshal(err)
    fmt.Fprintf(w, string(jsonError))
    return
  }
}
