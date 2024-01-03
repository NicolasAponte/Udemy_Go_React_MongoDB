package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func Registration(ctx context.Context) models.Response {
	var t models.User
	var r models.Response
	r.Status = 400

	fmt.Println("Entre a Registro")

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = err.Error()
		fmt.Println(r.Message)
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "Email is nul"
		fmt.Println(r.Message)
		return r
	}

	if len(t.Password) < 6 {
		r.Message = "Your password must have 6 characters at least"
		fmt.Println(r.Message)
		return r
	}

	_, exists, _ := bd.UserAlreadyExists(t.Email)
	if exists {
		r.Message = "User already exists"
		fmt.Println(r.Message)
		return r
	}

	_, status, err := bd.Registrate(t)
	if exists {
		r.Message = "Error registrating user " + err.Error()
		fmt.Println(r.Message)
		return r
	}
	if !status {
		r.Message = "It wasn't possible register the user"
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "Successful regitration"
	fmt.Println(r.Message)
	return models.Response{}
}
