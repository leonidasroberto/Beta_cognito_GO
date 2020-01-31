package main

import (
	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func (c CognitoConex) admAuth() {
	res, err := c.CognitoClient.AdminInitiateAuth(&cognito.AdminInitiateAuthInput{
		UserPoolId: aws.String(c.UserPoolID),
		ClientId:   aws.String(c.AppClientID),
	})

	if err != nil {
		println("Erro -> ", err)
		return
	}

	println(res)

}
