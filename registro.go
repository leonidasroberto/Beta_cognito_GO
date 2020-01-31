package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

///Função de registro
func (c CognitoConex) Registrar(w http.ResponseWriter, r *http.Request) {

	//Tratamento de requisições diferentes de POST
	if r.Method != "POST" {
		fmt.Fprint(w, "Ocorreu algum erro na sua solicitação!")
		return
	}
	var c_user cognitoUser

	c_user = cognitoUser{
		Email:    r.FormValue("user_email"),
		Password: r.FormValue("user_pass"),
		Username: r.FormValue("user_name"),
	}

	user := &cognito.SignUpInput{
		Username: aws.String(c_user.Username),
		Password: aws.String(c_user.Password),
		ClientId: aws.String(c.AppClientID),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(c_user.Email),
			},
		},
	}

	_, err := c.CognitoClient.SignUp(user)
	if err != nil {
		//fmt.Println("Erro aqui: ", err)
		fmt.Fprintln(w, err)
		return
	}

	//tpl.ExecuteTemplate(w, "index", user_name)
	redirec := "/confirma?user=" + c_user.Username

	http.Redirect(w, r, redirec, http.StatusSeeOther)
	return

}

///Função Confirmação de Email
func (c *CognitoConex) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	User := r.FormValue("user")
	if r.Method != "POST" {
		tpl.ExecuteTemplate(w, "confirm.html", User)
		return
	}
	code := r.FormValue("code")
	user_conf := r.FormValue("user_conf")

	dump := &cognito.ConfirmSignUpInput{
		Username:         aws.String(user_conf),
		ClientId:         aws.String("2va2f2bh1o36q4d2dt6f6qf7k9"),
		ConfirmationCode: aws.String(code),
	}

	_, err := c.CognitoClient.ConfirmSignUp(dump)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintln(w, "Usuário confirmado com sucesso!")

}
