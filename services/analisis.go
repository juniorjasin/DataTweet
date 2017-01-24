package services

import(
  "net/url"
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"
  "sort"

  "github.com/juniorjasin/datatweet/model"
)

/* Metodo donde comienza el proceso de calcular el porcentaje de favoritos
 * de una cuenta. Recive un url.Values donde obtengo los paraemtros de la url
 * que deben traer un Token y un Secret valido para poder crear un http.client
 * y poder realizar consultas a la API de twitter
 */
func PercentageOfFavorites(v url.Values) model.PairList {
  // creo un nuevo cliente
  client := model.GetClient(v)
  if client == nil {
    fmt.Println("\n\n ****CLIENTE NO INICIALIZADO**** \n\n")
    return nil
  }

  // tomo el screen_name que debe venir en la url para saber que cuenta voy a analizar
  sn := v.Get("screen_name")
  response, err := client.Get(
    // en count especifico la cantidad de tweets a analizar, puse el maximo (200)
		"https://api.twitter.com/1.1/favorites/list.json?count=200&screen_name="+ sn)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
  bits, err := ioutil.ReadAll(response.Body)

  // una vez obtenido el response lo 'pongo' en un puntero a interface[] para parsearlo
  var f []interface{} // tiene que ser un puntero a interface{} porque no es un solo JSON sino un array de JSON
  err1 := json.Unmarshal(bits, &f)
  if err1 != nil {
    fmt.Println(err1)
    return nil
  }

  // le paso f al metodo que devuelve una estructura PairList que es un puntero
  // de un struct string, float64 con los porcentajes y ordenada
  pl := calculatePercentage(f)

  return pl
}

//Calculo de porcentaje de favs
func calculatePercentage(f []interface{}) model.PairList {

  // Parseo el array de JSON y busco el array de 'user', donde especifica la informacion
  // propia del usuario al que se faveo, y obtengo su screen_name (su @).
  //  Creo un mapa donde sumo y obtengo screen_name:cantfavs y cuento la cantidad total en count
  cantNames := make(map[string]int)
  count := 0
  for _, v := range f {
    z := v.(map[string]interface{})
    for k2, v2 := range z {

      if k2 == "user" {
        count++
        u := v2.(map[string]interface{})
        s := u["screen_name"]
        sn := s.(string)
        j, ok := cantNames[sn]
        if ok == false {
          // si no existia, lo agrego y comienza con 1
          cantNames[sn] = 1
        }else{
          // si existia, j tiene el valor entonces agrego ++j
          j++
          cantNames[sn] = j
        }
      }
    }
  }

// recorro el mapa obtenido y creo un nuevo mapa con screen_name:porcentaje
  percentageMap := make(map[string]float64)
  for key, _ := range cantNames {
    i := cantNames[key]
    p := float64((i * 100))/ float64(count)
    fmt.Println("\n porcentaje:",p, " screen_name:", key)
    percentageMap[key] = p
  }

  // SortMap ordena el mapa y lo mete en el type PairList
  pl := SortMap(percentageMap)

  return pl
}

// ordena el mapa con los valores y devuelve un puntero a una estructura Pair con su clave y valor
func SortMap(wordFrequencies map[string]float64) model.PairList{
  pl := make(model.PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
    pl[i] = model.Pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}








//
