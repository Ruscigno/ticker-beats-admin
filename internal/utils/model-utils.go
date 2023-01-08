package utils

import (
	"reflect"
	"strings"
)

//StructToSlice return field name and its values from a struct
func StructToSlice(st interface{}, ignoreFields []string) ([]interface{}, []string) {
	var values []interface{}
	var fields []string
	v := reflect.Indirect(reflect.ValueOf(st))
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		fieldName := strings.ToLower(v.Type().Field(i).Name)
		toAdd := true
		for _, ig := range ignoreFields {
			if strings.EqualFold(fieldName, ig) {
				toAdd = false
				break
			}
		}
		if !toAdd {
			continue
		}
		fields = append(fields, fieldName)
		values = append(values, v.Field(i).Interface())
	}
	return values, fields
}
