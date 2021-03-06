package services

import(
  "sort"
)

// tipos y metodos para ordenar mapa
type Pair struct {
  Key string `json:"cadena"`
  Value float64 `json:"porcentaje"`
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

// ordena el mapa con los valores y devuelve un puntero a una estructura Pair con su clave y valor
func SortPairList(wordPercentage PairList) PairList{
  pl := make(PairList, len(wordPercentage))
  i := 0
  for _, v := range wordPercentage {
    pl[i] = Pair{v.Key, v.Value}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

// ordena el mapa con los valores y devuelve un puntero a una estructura Pair con su clave y valor
func SortMap(wordFrequencies map[string]float64) PairList{
  pl := make(PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
    pl[i] = Pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}
