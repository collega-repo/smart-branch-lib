package mapper

import (
	"github.com/goccy/go-json"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type Maps []Map

func (m Maps) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Maps) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type Map struct {
	Maps map[string]any `json:"maps"`
}

func (m Map) GetFloat(key string) float64 {
	v, _ := m.Maps[key].(float64)
	return v
}

func (m Map) GetInt(key string) int64 {
	switch v := m.Maps[key].(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case string:
		parseInt, _ := strconv.ParseInt(v, 10, 64)
		return parseInt
	default:
		return 0
	}
}

func (m Map) GetString(key string) string {
	v, _ := m.Maps[key].(string)
	return v
}

func (m Map) GetTimeStamp(key string) time.Time {
	var valueDateTime time.Time
	switch v := m.Maps[key].(type) {
	case time.Time:
		valueDateTime = v
	case string:
		valueDateTime, _ = time.Parse(`2006-01-02 15:04:05.000000`, v)
	}
	return valueDateTime
}

func (m Map) GetDate(key string) time.Time {
	var valueDate time.Time
	switch v := m.Maps[key].(type) {
	case time.Time:
		valueDate = v
	case string:
		valueDate, _ = time.Parse(`2006-01-02`, v[:10])
	}
	return valueDate
}

func (m Map) GetDecimal(key string) decimal.Decimal {
	var valueDecimal decimal.Decimal
	switch v := m.Maps[key].(type) {
	case string:
		valueDecimal, _ = decimal.NewFromString(v)
	case float64:
		valueDecimal = decimal.NewFromFloat(v)
	case int64:
		valueDecimal = decimal.NewFromInt(v)
	}
	return valueDecimal
}

func (m Map) GetBool(key string) bool {
	var valid bool
	switch v := m.Maps[key].(type) {
	case bool:
		valid = v
	case string:
		valid, _ = strconv.ParseBool(v)
	}
	return valid
}

func (m Map) Get(key string) any {
	return m.Maps[key]
}

func (m Map) GetMap(key string) Map {
	var dataMap Map
	a := m.Maps[key]
	dataMap.Maps, _ = a.(map[string]any)
	return dataMap
}

func (m Map) GetMaps(key string) (listMap []Map) {
	switch a := m.Maps[key].(type) {
	case []any:
		for i := range a {
			switch v := a[i].(type) {
			case map[string]any:
				listMap = append(listMap, Map{Maps: v})
			}
		}
	}
	return
}

func (m *Map) Put(key string, value any) {
	switch v := value.(type) {
	case Map:
		m.Maps[key] = v.Maps
	case map[string]any:
		m.Maps[key] = v
	default:
		m.Maps[key] = value
	}
}

func (m *Map) PutAll(maps Map) {
	m.PutAllMap(maps.Maps)
}

func (m *Map) PutAllMap(maps map[string]any) {
	for k, v := range maps {
		m.Maps[k] = v
	}
}

func (m Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Maps)
}

func (m *Map) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.Maps)
}

func (m *Map) Scan(value any) error {
	v, _ := value.([]byte)
	err := json.Unmarshal(v, &m.Maps)
	return err
}

func (m Map) MarshalBinary() ([]byte, error) {
	return json.Marshal(m.Maps)
}

func (m *Map) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m.Maps)
}
