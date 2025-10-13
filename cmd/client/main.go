package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api-2/internal/models"
)

func main() {
	req := models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	b, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:3000/users", "application/json", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var responseApi models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&responseApi); err != nil {
			panic(err)
		}
		panic(responseApi.Reason)
	}

	var responseApi models.CreateUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseApi); err != nil {
		panic(err)
	}

	fmt.Printf("New user ID: %v\n", responseApi.NewUserID)
}
