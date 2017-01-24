package controllers

import(
  "net/http"
  "fmt"
  "log"
  "encoding/json"
  //"net/url"

  "github.com/mrjones/oauth"
  "github.com/juniorjasin/datatweet/model"
  "github.com/juniorjasin/datatweet/services"
)
// variable auxiliar para no perder el valor de tokens
var tokensAux map[string]*oauth.RequestToken
var finalToken *oauth.AccessToken // lo guardo al pedo
var prueba string

/* Se llama a este metodo cuando llega un request a localhost:8888
 * metodo redirije al usuario para que acepte los permisos
 * siguiendo el estandar oauth1.0 de twitter. Utilizo la libreria "github.com/mrjones/oauth"
 * para realizarlo porque Go solo tiene el package oauth2 para oauth 2.0
 */
func RedirectUserToTwitter(w http.ResponseWriter, r *http.Request) {
  var c *oauth.Consumer
  var tokens map[string]*oauth.RequestToken
  c, tokens = model.GetConsumer()
  tokensAux = tokens
	tokenUrl := fmt.Sprintf("http://%s/maketoken", r.Host)
	token, requestUrl, err := c.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
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
  c, tokens = model.GetConsumer()
  tokens = tokensAux

	values := r.URL.Query()
	verificationCode := values.Get("oauth_verifier")
	tokenKey := values.Get("oauth_token")

	accessToken, err := c.AuthorizeToken(tokens[tokenKey], verificationCode)
	if err != nil {
		log.Fatal(err)
	}

  client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
  }

  cli := model.SetClient(client)
  if cli == nil{
    log.Fatal(cli)
  }

  // guardo el token con el que voy a hacer consultas
  finalToken = accessToken

  // *************** CREO QUE NO NECESITO ESTO SI ES QUE YA TENGO EL CLIENTE *********************
  bits2, err := json.Marshal(accessToken)
  if err != nil{
    log.Fatal(err)
  }
  // Con este metodo visualizo la respuesta en el navegador, retorna el json con el token
  fmt.Fprintf(w, string(bits2))
}


/* Parseo los parametros de la url se los paso a la capa de services
   Tengo que saber como van a ser los json que va a recibir para poder extraer los values


      "Token":"811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16",
      "Secret":"p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv",
      "screen_name":"juniorjasin",
      "user_id":"811672150000209920",

      "Token":"2439545395-84074Ec8VsIXygKaUdt58mobvBNvE8A4Bu5ZTZq",
      "Secret":"qP837ZmljgyJoRVWJyRS1ujzf2EwR3x2BOxeIeLyLJJmb",
      "AdditionalData":{"screen_name":"EsJorgito","user_id":"2439545395","x_auth_expires":"0"}}
*/

/*
 * obtuve un mapa de los parametros que van a venir en la url,
 * ahora se los paso a la capa de services para que le pegue a la api
 * de twitter y asi obtener datos y procesarlos y devolver una estructura
 * que tengo que hacer y que debe estar implementada en la capa de model para
 * estructurar los datos (podria ser % favs para cada persona por ej) y
 * desde aca (capa controllers) devolver la respuesta a la app que consuma esta api
*/

/* url_mappings llama a este metodo cuando llega un request a /GetPercentageOfFavorites
 * Primero se obtienen los paraemtros de la url y se los paso a la capa de services.
 * Luego Respondo los primeros 10 resultados que ya vienen ordenados
 */
func GetPercentageOfFavorites(w http.ResponseWriter, r *http.Request){
  values := r.URL.Query()
  pl := services.PercentageOfFavorites(values)
  if pl == nil {
    fmt.Fprintf(w, "ERROR FALTAN PARAMETROS: Token y/o Secret...") // nombre
  }

  i := 0
  for _, value := range pl {
      if i > 9 { break }
      fmt.Fprintf(w, value.Key +" ") // nombre
      fmt.Fprintf(w, "%.2f", value.Value) // porcentaje
      fmt.Fprintf(w,"%v", "%")
      fmt.Fprintf(w, "\n")
      i++
  }
}

/* DATOS IVAN:

{"Token":"635962800-h2aiHePot0uqEJzals0zPuSnLV0S6y2mvIRuoAKI",
"Secret":"72WCC1YKW2pXOZDdNNgw1U4aYJMqLRWXxPccwoNKp9JHR",
"AdditionalData":{"screen_name":"Ds_Ivan_Gs","user_id":"635962800","x_auth_expires":"0"}}

DATOS juniorjasin
/*
"Token":"811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16",
"Secret":"p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv",
"screen_name":"juniorjasin",
"user_id":"811672150000209920",
*/
