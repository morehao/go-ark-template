package utils

import (
	"encoding/json"
	"fmt"
)

func MapToJson() {
	student := make(map[string]interface{})
	student["name"] = "5lmh.com"
	student["age"] = 18
	student["sex"] = "man"
	b, err := json.Marshal(student)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}

func JsonToStruct() {
	type Person struct {
		Age       int    `json:"age,string"`
		Name      string `json:"name"`
		Niubility bool   `json:"niubility"`
	}
	b := []byte(`{"age":"18","name":"5lmh.com","marry":false}`)
	var p Person
	err := json.Unmarshal(b, &p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}
