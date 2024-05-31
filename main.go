package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/juanmontilva/xgo/awsgo"
	"github.com/juanmontilva/xgo/bd"
	"github.com/juanmontilva/xgo/handlers"
	"github.com/juanmontilva/xgo/models"
	"github.com/juanmontilva/xgo/secretmanager"
)

func main() {
	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	if !ValidoParametro() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno. debe incluir 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret" + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	path := strings.Replace(request.PathParameters["xclone"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtsign"), SecretModel.JWTSign)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	//Chequeo conexion a la base de datos

	err = bd.ConectarBD(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error al conectar en la base de datos" + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	//res api
	resAPI := handlers.Manejadores(awsgo.Ctx, request)
	if resAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: resAPI.Status,
			Body:       resAPI.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	} else {
		return resAPI.CustomResp, nil
	}
}

func ValidoParametro() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}

	return traeParametro

}
