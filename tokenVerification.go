package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

////Verificador de token
func (c CognitoConex) verificAuth(token string) {

	res, err := c.CognitoClient.VerifySoftwareToken(&cognito.VerifySoftwareTokenInput{
		AccessToken: aws.String(token)})

	if err != nil {
		println("ERRO DE AUTENTICAÇÃO!")
		println(err)
		return
	}

	println("TOKEN VÀLIDO!")
	println(res)

}

///Função de autenticação
func (c CognitoConex) Aut(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Fprintln(w, "Ocorreu algum erro na sua requisição!")
		return
	}

	username := r.FormValue("user")
	password := r.FormValue("pass")
	/*
		println("Usuario: " + username)
		println("Senha: " + password)
		println("CliendIG: " + c.AppClientID)
	*/
	flow := aws.String(att)

	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(c.AppClientID),
	}

	res, err := c.CognitoClient.InitiateAuth(authTry)
	if err != nil {
		fmt.Fprintln(w, err)
		//http.Redirect(w, r, fmt.Sprintf("/login?error=%s", err.Error()), http.StatusSeeOther)
		return
	}

	var returnJSON infoAuth
	returnJSON = tratamentoAuth(res.String())
	fmt.Fprintln(w, "Token de acesso! -> "+returnJSON.AccessToken)
	fmt.Fprintln(w, "Retorno id -> "+returnJSON.IdToken)

}
