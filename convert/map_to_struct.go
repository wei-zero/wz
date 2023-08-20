package convert

import "encoding/json"

// MapToStruct converts a map to a struct.
func MapToStruct(m map[string]interface{}, s interface{}) error {
	b, _ := json.Marshal(m)
	return json.Unmarshal(b, s)
}

func AnyToStruct(m any, s interface{}) error {
	b, _ := json.Marshal(m)
	return json.Unmarshal(b, s)
}

func StructToMap(s interface{}, m map[string]interface{}) error {
	b, _ := json.Marshal(s)
	return json.Unmarshal(b, &m)
}
