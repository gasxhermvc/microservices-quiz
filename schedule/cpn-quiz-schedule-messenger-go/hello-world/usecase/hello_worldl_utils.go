package usecase

import (
	"reflect"
)

func (hello helloWorldUseCase) structToMap(data any) map[string]interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := reflect.TypeOf(data)

	result := make(map[string]interface{})
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		value := dataValue.Field(i).Interface()
		tag := field.Tag.Get("struct")

		if tag != "" {
			result[tag] = value
		}
	}

	return result
}
