package utils

import "reflect"

func MakeObjectWithoutNilFields(modelData interface{}) map[string]interface{} {
	data := make(map[string]interface{})

	// Converte para reflexão
	v := reflect.ValueOf(modelData)
	t := reflect.TypeOf(modelData)

	// Se for um ponteiro, pega o valor que ele aponta
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Se o campo for um ponteiro, verifica se é nil
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue // Ignora campos nil
			}
			data[field.Name] = value.Elem().Interface()
		} else {
			data[field.Name] = value.Interface()
		}
	}

	return data
}
