package appdataservice

import (
	"cpn-quiz-api-authentication-go/utils"
	"fmt"
)

type QueryResult struct {
	Total           int
	OutputParameter map[string]interface{}
	Message         string
	Success         bool
	Error           error
}

func (queryResult QueryResult) ConstainsKey(parameterName string) bool {
	if _, ok := queryResult.OutputParameter[parameterName]; ok {
		return true
	}
	return false
}

func (queryResult QueryResult) AddOutputParameter(parameterName string, value interface{}) {
	if !queryResult.ConstainsKey(parameterName) {
		queryResult.OutputParameter[parameterName] = value
	} else {
		fmt.Println(fmt.Sprintf("ไม่สามารถเพิ่ม Output Parameter \"%s\" ได้เนื่องจากมีอยู่แล้ว", parameterName))
	}
}

func (queryResult QueryResult) RemoveOutputParameter(parameterName string) {
	if queryResult.ConstainsKey(parameterName) {
		delete(queryResult.OutputParameter, parameterName)
	} else {
		fmt.Println(fmt.Sprintf("ไม่สามารถลบได้, เนื่องจากไม่พบ Output Parameter \"%s\"", parameterName))
	}
}

func (queryResult QueryResult) UpdateOutputParameter(parameterName string, value interface{}) {
	if queryResult.ConstainsKey(parameterName) {
		queryResult.OutputParameter[parameterName] = value
	} else {
		fmt.Println(fmt.Sprintf("ไม่สามารถลบได้, เนื่องจากไม่พบ Output Parameter \"%s\"", parameterName))
	}
}

func (queryResult QueryResult) GetOutputParameter(parameterName string) interface{} {
	if queryResult.ConstainsKey(parameterName) {
		return queryResult.OutputParameter[parameterName]
	}

	return nil
}

func (queryResult QueryResult) ToJson() (string, error) {
	return utils.JsonSerialize(queryResult)
}

func (queryResult QueryResult) ToBytes() ([]byte, error) {
	json, err := queryResult.ToJson()
	if err != nil {
		fmt.Println(err.Error())
	}
	bytes := []byte(json)
	return bytes, nil
}

func (queryResult QueryResult) ToDictionary() map[string]interface{} {
	return queryResult.OutputParameter
}

func (queryResult QueryResult) FieldRename(columns []string) map[string]interface{} {
	results := make(map[string]interface{})

	idx := 0
	for key, dataset := range queryResult.OutputParameter {
		if columns[idx] != "" {
			results[columns[idx]] = dataset
		} else {
			results[key] = dataset
		}
		idx++
	}

	return results
}
