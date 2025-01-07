package helper

import (
	"encoding/json"
	"maze-conquest-api/model/domain"
)

func NewUser(data map[string]interface{}) *domain.User {
	// 01 Convert your map to JSON and then unmarshal it into the struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// 02 Declare the variable; manually assign for UUID
	var user domain.User

	// 03 Automatically assign to the properties with json.Unmarshal
	err = json.Unmarshal(jsonBytes, &user)
	if err != nil {
		panic(err)
	}

	return &user
}

// func NewHero() {}
