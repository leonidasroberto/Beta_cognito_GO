package main

import (
	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func (c CognitoConex) deletUser(tokenAccess string) {
	del := &cognito.DeleteUserInput{
		AccessToken: aws.String(tokenAccess),
	}

	ret, err := c.CognitoClient.DeleteUser(del)
	if err != nil {
		println("Erro -> ", err)
	}

	println(ret)

}
