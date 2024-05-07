package main

import (
	"fmt"
	"strings"

	group "github.com/Drelf2018/gin-group"
	"github.com/gin-gonic/gin"
)

func GetPing(ctx *gin.Context) (data any, err group.Error) {
	return "pong", nil
}

type User string

var ErrInvalidUser = group.NewError(1, "%s is an invalid username")

func SetUser(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "Bob" {
		group.Abort(ctx, nil, ErrInvalidUser.Format(name))
		return
	}
	group.SetUser(ctx, User(strings.TrimSpace(name)))
}

func GetHello(ctx *gin.Context) (any, group.Error) {
	user := group.GetUser[User](ctx)
	if user == "" {
		user = "Anonymous"
	}
	return fmt.Sprintf("Hello %s!", user), nil
}

func main() {
	user := group.Group{
		Path: "user",
		Middlewares: gin.HandlersChain{
			SetUser,
		},
		Handlers: group.Chain{
			GetHello,
		},
	}

	visitor := group.Group{
		Middlewares: gin.HandlersChain{
			group.CORS,
			group.Static("vue"),
		},
		Handlers: group.Chain{
			GetPing,
			GetDownload,
			GetResourceColonfile,
			GetPicColonmid,
		},
		Groups: []group.Group{user},
	}

	group.Default(visitor).Run("localhost:8080")
}
