package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/juanmontilva/xgo/jwt"
	"github.com/juanmontilva/xgo/models"
	"github.com/juanmontilva/xgo/routers"
)

func Manejadores(ctx context.Context, request events.APIGatewayProxyRequest) models.RespApi {
	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + ">" + ctx.Value(models.Key("method")).(string))

	var r models.RespApi

	r.Status = 400

	isOk, statusCode, msg, _ := validoAuthorization(ctx, request)

	if !isOk {
		r.Status = statusCode
		r.Message = msg
		return r
	}

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "registro":
			return routers.Registro(ctx)

		case "login":
			return routers.Login(ctx)
		}
		//
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {

		}

		//
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {

		}
		//
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {

		}

		//

	}

	r.Message = "Method Invalid"
	return r

}

func validoAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	path := ctx.Value(models.Key("path")).(string)
	if path == "registro" || path == "login" || path == "obtenerAvatar" || path == "obtenerBanner" {
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, 401, "TOKEN REQUERIDO", models.Claim{}
	}

	claim, todoOK, msg, err := jwt.ProcesoToken(token, ctx.Value(models.Key("jwtSign")).(string))

	if !todoOK {
		if err != nil {
			fmt.Println("ERROR EN EL TOKEN" + err.Error())
			return false, 401, err.Error(), models.Claim{}
		} else {
			fmt.Println("ERROR EN EL TOKEN" + msg)
			return false, 401, msg, models.Claim{}
		}
	}

	fmt.Println("TOKEN OK")

	return true, 200, msg, *claim

}
