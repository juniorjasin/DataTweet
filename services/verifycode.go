package services

import(
  "encoding/json"
)

// verifico el codigo de error que me devuelve la API de twitter
func verifyCode(bits []byte) float64 {
  var f interface{}
  json.Unmarshal(bits, &f)

  m := f.(map[string]interface{})
  for _,v := range m {
    val := v.([]interface{})
    for _,v2 := range val {
      m2 := v2.(map[string]interface{})
      c := m2["code"]
      code := c.(float64)
      return code
    }
  }

  return 0
}
