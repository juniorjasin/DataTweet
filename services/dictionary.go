package services

import(
  "fmt"
  "log"
  "strings"
  "unicode"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "golang.org/x/text/transform"
  "golang.org/x/text/unicode/norm"
)

var greaterId float64 = 0
var errorCode float64

// metodo principal que retorna el diccionario con palabras y porcentaje
func GetDictionary(client *http.Client, sn string) (PairList, float64) {

  bits := getResponseDic(client, sn)
  json := getJsonDic(bits)
  if json == nil{
    return nil,verifyCode(bits)
  }
  tweets := getTweets(json)
  dic := getMostUsedWords(tweets)
  return dic, 0
}

// obtengo el response segun el valor de si
func getResponseDic(client *http.Client, sn string) []byte{
  var bits []byte

  response, err := client.Get(
      // en count especifico la cantidad de tweets a analizar, puse el maximo (200)
  	"https://api.twitter.com/1.1/statuses/user_timeline.json?include_rts=false&count=20&screen_name="+ sn)
  if err != nil {
  	log.Fatal(err)
  }
  defer response.Body.Close()
  bits, _ = ioutil.ReadAll(response.Body)

  return bits
}

// Hago unamarshal del json y verifico que sea correcto
func getJsonDic(bits []byte) []interface {}{
  // una vez obtenido el response lo 'pongo' en un puntero a interface[] para parsearlo
  var f []interface{} // tiene que ser un puntero a interface{} porque no es un solo JSON sino un array de JSON
  err1 := json.Unmarshal(bits, &f)
  if err1 != nil {
    fmt.Println(err1)
    return nil
  }

  return f
}


// obtengo un array con cada tweet, filtrando los RT, Menciones y las que contengan links, fotos, gif, etc.
func getTweets(f []interface{}) []string {

  tweet := make([]string, 0)
  first := true
  for _, v := range f {
    z := v.(map[string]interface{})
    for k2, v2 := range z {
      if first == true && k2 == "id"{
        greaterId = v2.(float64) // podria retornar esta variable para que la mande el usuario de nuevo
        first = false
      }

      if k2 == "text" {
        str := v2.(string)
        // filtro los retweets y las menciones
        if !strings.Contains(str, "RT") && !strings.Contains(str, "http") && !strings.Contains(str, "@"){
          tweet = append(tweet, str)
        }
      }
    }
  }

  return tweet
}

// obtengo las 10 palabras mas utilizadas del array de palabras que viene
func getMostUsedWords(tweets []string) PairList {
  words := getWords(tweets)
  mapWords, count := getWordsFrequency(words)
  dic := getWordsPercentage(mapWords, count)
  pl := SortPairList(dic)
  pl = pl[:10]

  return pl
}

// obtengo todas las palabras
func getWords(tw []string) []string{
  words := make([]string, 0)
  for _,v := range tw {
    str := strings.Fields(v)
    for _,s := range str {
      words = append(words, s)
    }
  }

  words = removeAccent(words)
  words = toLowerCase(words)
  words = filterWords(words)

  return words
}

// metodo que remueve todos los acentos que encuentre en cada string del slice de string
func removeAccent(words []string)[]string {
  for i,s := range words {
    t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
    result, _, _ := transform.String(t, s)
    words[i] = result
  }

  return words
}

// metodo auxiliar para remover los acentos
func isMn(r rune) bool {
    return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// convierto todos los strings del slice a minuscula
func toLowerCase(words []string)[]string  {
  for i,s := range words {
    words[i] = strings.ToLower(s)
  }
  return words
}

// var innecesaryWords []string = []string{"tu", "el", "nosotros", "yo"}
// metodo para filtrar pronombes y articulos en las palabras que obtuve
func filterWords(words []string) []string {
  innecesaryWords := []string{"tu", "el", "nosotros", "yo", "la", "los", "las", "de", "es"} //
  for i := 0; i < len(words); i++ {
    for _,s := range innecesaryWords {
      if words[i] == s {
        words = append(words[:i], words[i+1:]...)
        i--
      }
    }
  }
  return words
}

// obtengo mapa de [palabras]:apariciones en tweets
func getWordsFrequency(words []string) (map[string]int, int) {
  mapWords := make(map[string]int)
  count := 0
  for _,v := range words {
    // fmt.Println("word [",i,"]:", v)
    c, ok := mapWords[v]
    if ok == true{
      c++
      mapWords[v] = c
    } else {
      mapWords[v] = 1
    }
    count++
  }
  return mapWords,count
}

// calculo mapa [palabra]:porcentaje de apariciones
func getWordsPercentage(mapWords map[string]int, count int) PairList{
  dic := make(PairList, 0)
  for word, cant := range mapWords {
    // fmt.Println("word:", word, " cant:", cant)
    per := float64((cant * 100)) / float64(count)
    dic = append(dic, Pair{word, per})
  }
  return dic
}
