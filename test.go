package main

import "fmt"

type Response struct {
	Status  int
	Message string
	Data    DuplicateResponse
}

type DuplicateResponse struct {
	Token        string
	RefreshToken string
}

func main() {
	var Re Response
	var Du DuplicateResponse
	Du.Token = "good"
	Du.RefreshToken = "very good"

	Re.Status = 200
	Re.Message = "qwe"
	Re.Data = Du

	fmt.Println(Re)
}
