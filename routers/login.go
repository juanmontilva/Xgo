package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/juanmontilva/xgo/bd"
	"github.com/juanmontilva/xgo/jwt"
	"github.com/juanmontilva/xgo/models"
)

func Login(ctx context.Context) models.RespApi {
	var t models.Usuario
	var r models.RespApi
	r.Status = 400

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "USUARIO Y CONTRASEÑA NO VALIDO" + err.Error()
		return r
	}
	if len(t.Email) == 0 {
		r.Message = "EL EMAIL DEL USUARIO ES REQUERIDO"
		return r
	}

	userData, existe := bd.IntentoLogin(t.Email, t.Password)
	if !existe {
		r.Message = "USUARIO Y/O CONSTRASEÑA INVALIDOS"
		return r
	}

	jwtKey, err := jwt.GeneroJWT(ctx, userData)
	if err != nil {
		r.Message = "ocurrio un erro ral generar token correspondiente > " + err.Error()
		return r
	}

	resp := models.RespuestaLogin{
		Token: jwtKey,
	}

	token, err2 := json.Marshal(resp)

	if err2 != nil {
		r.Message = "OCURRIO UN ERROR AL INTENTAR FORMATEAR EL TIKEN A JSON > " + err2.Error()
		return r
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}

	cookieString := cookie.String()

	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieString,
		},
	}

	r.Status = 200
	r.Message = string(token)
	r.CustomResp = res

	return r

}
