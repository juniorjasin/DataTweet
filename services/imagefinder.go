package services

import(
  "log"
  "fmt"
  "net/url"
  "net/http"
  "io/ioutil"
  "encoding/json"
  // "reflect"

  "github.com/juniorjasin/datatweet/model"
)

func GetImages(client *http.Client, v url.Values) model.UrlNameList {

  fmt.Println("\n\n\n GETIMAGES: \n")
  // tomo el screen_name que debe venir en la url para saber que cuenta voy a analizar
  // sn := v.Get("screen_name")
  response, err := client.Get(
		"https://api.twitter.com/1.1/statuses/home_timeline.json?count=50")
	if err != nil {
		log.Fatal(err)
    fmt.Println("\n\n\n HUBO UN ERROR \n\n")
	}
	defer response.Body.Close()
  bits, err := ioutil.ReadAll(response.Body)

  fmt.Println("\n\n\n RESPUESTA: \n")
  fmt.Println(string(bits))

  // una vez obtenido el response lo 'pongo' en un puntero a interface[] para parsearlo
  var f []interface{} // tiene que ser un puntero a interface{} porque no es un solo JSON sino un array de JSON
  err1 := json.Unmarshal(bits, &f)
  if err1 != nil {
    fmt.Println(err1)
  }

  return getUrls(f)
}

/* Recibo un puntero a interface{} donde esta el json del response.
 * Parseo el json y extraigo las urls de cada imagen
 */
func getUrls(f []interface{}) model.UrlNameList {

  // 3200 porque es el max de tweets que puedo consultar en el timeline,
  // pero podria ser dinamico y asignarse solo lo que vine por parametros en la url
  snames := make([]string, 3200)
  count1 := 0
  count := 0


// primero obtengo un mapa determinando la [index]:screen_name
  for _, v := range f {
    z := v.(map[string]interface{})
    for k2, v2 := range z {
      if k2 == "user" {
        u := v2.(map[string]interface{})
        s := u["screen_name"]
        sn := s.(string)
        snames[count1] = sn
        count1++
      }
    }
   }

urlnamelist := make(model.UrlNameList, count1)

// obtengo las url de las fotos las meto
  for _, v := range f {
    z := v.(map[string]interface{})
    for k2, v2 := range z {
      if k2 == "entities" {
        sndLevel := v2.(map[string]interface{})
        for k3,v3 := range sndLevel {
            if k3 == "media"{
            media := v3.([]interface{})
             for _,v4 := range media {
                urls := v4.(map[string]interface{})
                for k5,v5 := range urls{
                  if k5 == "media_url"{
                    fmt.Println("media_url:", v5.(string), "\n")
                    urlnamelist[count] = model.UrlName{count, v5.(string), snames[count]}
                    count++
                  }
                }
              }
            }
        }
      }
    }
  }

  urlnamelist = urlnamelist[:count]
  fmt.Println("COMIENZA MAPA: \n")
  for _, urlname := range urlnamelist{
    fmt.Println("idx:", urlname.Idx, "url:", urlname.Url, "snames:", urlname.Screen_name)
  }

  return urlnamelist
}
