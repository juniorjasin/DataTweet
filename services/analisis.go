package services

import(
  "net/url"
  "fmt"
  "log"
  "io/ioutil"
  "encoding/json"
  // "reflect"
  "sort"
  // "strings"
  // "io"

  "github.com/juniorjasin/datatweet/model"
)

/*
"Token":"811672150000209920-fTCkCDAbXD9NykbRY9NheMENYHJNA16",
"Secret":"p73tL8y3RJFchHqwn9uwsRJD34NPWkiBHxX3G3q0VE1zv",
"screen_name":"juniorjasin",
"user_id":"811672150000209920",
*/


func PercentageOfFavorites(v url.Values) model.PairList {

  // ******************************* importante **************************************
  // en GetClient debo pasarle los parametros Token y Secret para crear un nuevo cliente
  // dentro de GetClient voy a usar GetConsumer tambien y asi crear un consumer que llame a MakeHttpClient
  // y ahi creo un client y lo devuelvo, por ahora eso no esta hecho.
  client := model.GetClient()
  if client == nil{
    fmt.Println("\n\n ****CLIENTE NO INICIALIZADO**** \n\n")
  }

  sn := v.Get("screen_name")
  response, err := client.Get(
		"https://api.twitter.com/1.1/favorites/list.json?count=200&screen_name="+ sn)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
  bits, err := ioutil.ReadAll(response.Body)

  pl := calculatePercentage(bits)

  return pl
}

func calculatePercentage(b []byte) model.PairList {
  //*** podria pasar esto al metodo que lo llama y aca solo calcular
  var f []interface{} // tiene que ser un puntero a interface{} porque no es un solo JSON sino un array de JSON
  err := json.Unmarshal(b, &f)
  if err != nil {
    fmt.Println(err)
    return nil
  }
  //***
  fmt.Println("\n\n\n")

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
  fmt.Println("\n\n\n")
  fmt.Println("MAPA:", cantNames)

  fmt.Println("\n\n\n")
  fmt.Println("count:", count)


  percentageMap := make(map[string]float64)
  for key, _ := range cantNames {
    i := cantNames[key]
    p := float64((i * 100))/ float64(count)
    fmt.Println("\n porcentaje:",p, " screen_name:", key)
    percentageMap[key] = p
  }

  fmt.Println("\n\n\n")

  pl := rankByWordCount(percentageMap)
  // muestro el mapa de porcentaje solo para comprobar
  for _, value := range pl {
      fmt.Println("key:", value.Key) // nombre
      fmt.Println("value:",value.Value) // porcentaje
      fmt.Println(" ")
  }

  return pl
}

// ordena el mapa con los valores y devuelve un puntero a una estructura Pair con su clave y valor
func rankByWordCount(wordFrequencies map[string]float64) model.PairList{
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
