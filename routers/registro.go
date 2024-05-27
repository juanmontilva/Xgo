package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/juanmontilva/xgo/bd"
	"github.com/juanmontilva/xgo/models"
)

func Registro(ctx context.Context) models.RespApi {
	var t models.Usuario
	var r models.RespApi
	r.Status = 400

	fmt.Println("Entre a registro")

	body := ctx.Value(models.Key("body")).(string)

	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = err.Error()
		fmt.Println(r.Message)
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "DEBE ESPECIFICAR EL EMAIL"
		fmt.Println(r.Message)
		return r
	}

	if len(t.Password) < 6 {
		r.Message = "DEBE ESPECIFICAR UNA CONTRASEÑA DE AL MENOS 6 CARACTERES"
		fmt.Println(r.Message)
		return r
	}

	_, encontrado, _ := bd.ChequeoYaExisteUsuario(t.Email)
	if encontrado {
		r.Message = "EXISTE USUARIO REGISTRADO EN EMAIL"
		fmt.Println(r.Message)
		return r
	}

	_, status, err := bd.InsertoRegistro(t)
	if err != nil {
		r.Message = "OCURRIO UN ERROR AL INTENTAR REALIZAR EL REGISTRO DE USUARIO" + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "NO SE HA LOGRADO INSERTAR EL REGISTRO DEL USUARIO"
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "REGISTRO OK"
	fmt.Println(r.Message)
	return r

}
