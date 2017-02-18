package services

import(
  "fmt"
  "log"
  "net/url"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

/* Metodo donde comienza el proceso de calcular el porcentaje de favoritos
 * de una cuenta. Recive un url.Values donde obtengo los paraemtros de la url
 * que deben traer un Token y un Secret valido para poder crear un http.client
 * y poder realizar consultas a la API de twitter
 */
func PercentageOfFavorites(v url.Values, client *http.Client) (PairList,float64) {

  // tomo el screen_name que debe venir en la url para saber que cuenta voy a analizar
  sn := v.Get("screen_name")
  bits := getResponseFav(client, sn)
  json := getJsonFav(bits)
  if json == nil{
    return nil, verifyCode(bits)
  }
  pl := calculatePercentage(json)

  return pl, 0
}

func getResponseFav(client *http.Client, sn string) []byte {
  response, err := client.Get(
    // en count especifico la cantidad de tweets a analizar, puse el maximo (200)
		"https://api.twitter.com/1.1/favorites/list.json?count=200&screen_name="+ sn)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
  bits, err := ioutil.ReadAll(response.Body)
  return bits
}

// Hago unamarshal del json y verifico que sea correcto
func getJsonFav(bits []byte) []interface {}{
  // una vez obtenido el response lo 'pongo' en un puntero a interface[] para parsearlo
  var f []interface{} // tiene que ser un puntero a interface{} porque no es un solo JSON sino un array de JSON
  err1 := json.Unmarshal(bits, &f)
  if err1 != nil {
    fmt.Println(err1)
    return nil
  }
  return f
}

//Calculo de porcentaje de favs
func calculatePercentage(f []interface{}) PairList {
  nameCount, count := getNameCount(f)
  percentageMap := getPercentageMap(nameCount, count)
  pl := SortMap(percentageMap)

  return pl
}

// Parseo el array de JSON y busco el array de 'user', donde especifica la informacion
// propia del usuario al que se faveo, y obtengo su screen_name (su @).
//  Creo un mapa donde sumo y obtengo screen_name:cantfavs y cuento la cantidad total en count
func getNameCount(f []interface{}) (map[string]int, int) {
  nameCount := make(map[string]int)
  count := 0
  for _, v := range f {
    z := v.(map[string]interface{})
    for k2, v2 := range z {

      if k2 == "user" {
        count++
        u := v2.(map[string]interface{})
        s := u["screen_name"]
        sn := s.(string)
        j, ok := nameCount[sn]
        if ok == false {
          // si no existia, lo agrego y comienza con 1
          nameCount[sn] = 1
        }else{
          // si existia, j tiene el valor entonces agrego ++j
          j++
          nameCount[sn] = j
        }
      }
    }
  }
  return nameCount, count
}

// recorro el mapa obtenido y creo un nuevo mapa con screen_name:porcentaje
func getPercentageMap(nc map[string]int, count int) map[string]float64 {
  percentageMap := make(map[string]float64)
  for key, _ := range nc {
    i := nc[key]
    p := float64((i * 100))/ float64(count)
    // fmt.Println("\n porcentaje:",p, " screen_name:", key)
    percentageMap[key] = p
  }
  return percentageMap
}
