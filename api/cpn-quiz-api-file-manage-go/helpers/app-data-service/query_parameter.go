package appdataservice

import "fmt"

type QueryParameter struct {
	Parameters map[string]interface{}
}

func NewQueryParameter(parameters ...map[string]interface{}) QueryParameter {
	_parameters := make(map[string]interface{})

	for _, m := range parameters {
		for k, v := range m {
			_parameters[k] = v
		}
	}

	return QueryParameter{
		Parameters: _parameters,
	}
}

func (queryParameter QueryParameter) ConstainsKey(parameterName string) bool {
	if _, ok := queryParameter.Parameters[parameterName]; ok {
		return true
	}
	return false
}

func (queryParameter QueryParameter) Add(parameterName string, value interface{}) {
	if !queryParameter.ConstainsKey(parameterName) {
		queryParameter.Parameters[parameterName] = value
	} else {
		fmt.Println(fmt.Sprintf("ไม่สามารถเพิ่ม Query Parameter \"%s\" ได้เนื่องจากมีอยู่แล้ว", parameterName))
	}
}

func (queryParameter QueryParameter) Update(parameterName string, value interface{}) {
	if queryParameter.ConstainsKey(parameterName) {
		queryParameter.Parameters[parameterName] = value
	} else {
		fmt.Println(fmt.Sprintf("ไม่สามารถลบได้, เนื่องจากไม่พบ Query Parameter \"%s\"", parameterName))
	}
}

func (queryParameter QueryParameter) Remove(parameterName string) {
	if queryParameter.ConstainsKey(parameterName) {
		delete(queryParameter.Parameters, parameterName)
	} else {
		fmt.Println(fmt.Sprintf("ไม่สามารถลบได้, เนื่องจากไม่พบ Query Parameter \"%s\"", parameterName))
	}
}

func (queryParameter QueryParameter) RemoveAll() {
	for k := range queryParameter.Parameters {
		delete(queryParameter.Parameters, k)
	}
}
