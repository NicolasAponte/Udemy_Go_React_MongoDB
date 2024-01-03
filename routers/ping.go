package routers

import (
	"context"

	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

func Ping(ctx context.Context) models.Response {
	return models.Response{
		Status:  200,
		Message: "pong",
	}
}
