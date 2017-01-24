package controllers

import(
  "net/http"
  "fmt"
  "log"
  "reflect"
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

// metodo que redirije al usuario para que acepte los permisos
// siguiente el estandar oauth1.0 de twitter
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

// obtengo el accessToken luego de que se acceptaran los permisos
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

  fmt.Println("TIPO DE AccessToken:")
  fmt.Println(reflect.TypeOf(accessToken))

  client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
  }

  fmt.Println("TIPO DE CLIENTE")
  fmt.Println(reflect.TypeOf(client))

  cli := model.SetClient(client)
  if cli == nil{
    log.Fatal(cli)
  }

  // guardo el token con el que voy a hacer consultas
  finalToken = accessToken

/*
	response, err := client.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json?count=1")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()


	bits, err := ioutil.ReadAll(response.Body)
  fmt.Printf("\ntipo de variable bits:")
  reflect.TypeOf(bits)
  //*/

//  fmt.Fprintf(w, accessToken.Token) // aca respongo el accessToken

  // *************** CREO QUE NO NECESITO ESTO SI ES QUE YA TENGO EL CLIENTE *********************
  bits2, err := json.Marshal(accessToken)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Println("TIPO DEL ACCESSTOKEN")
  fmt.Println(reflect.TypeOf(accessToken))
  // retorno el json. Con este metodo veo la respuesta en el navegador
  fmt.Fprintf(w, "json de respuesta:" + string(bits2))
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

// obtuve un mapa de los parametros que van a venir en la url,
// ahora se los paso a la capa de services para que le pegue a la api
// de twitter y asi obtener datos y procesarlos y devolver una estructura
// que tengo que hacer y que debe estar implementada en la capa de model para
// estructurar los datos (podria ser % favs para cada persona por ej) y
// desde aca (capa controllers) devolver la respuesta a la app que consuma esta api
func GetPercentageOfFavorites(w http.ResponseWriter, r *http.Request){
  values := r.URL.Query()
  pl := services.PercentageOfFavorites(values)

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













//
