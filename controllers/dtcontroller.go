package controllers

import(
  "net/http"
  "fmt"
  "log"
  "io/ioutil"
  "reflect"
  "encoding/json"

  "github.com/mrjones/oauth"
  "github.com/juniorjasin/datatweet/model"
)
// variable auxiliar para no perder el valor de tokens
var tokensAux map[string]*oauth.RequestToken

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

	client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json?count=1")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
  fmt.Printf("\ntipo de variable bits:")
  reflect.TypeOf(bits)
//  fmt.Fprintf(w, accessToken.Token) // aca respongo el accessToken
  bits2, err := json.Marshal(accessToken)
  if err != nil{
    log.Fatal(err)
  }

  // retorno el json. Con este metodo veo la respuesta en el navegador
  fmt.Fprintf(w, "json de respuesta:" + string(bits2))
}
