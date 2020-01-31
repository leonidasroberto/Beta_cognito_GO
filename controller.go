package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type cognitoUser struct {
	Username string
	Password string
	Email    string
}

///Infomações de login
type infoAuth struct {
	AccessToken  string
	ExpiresIn    string
	IdToken      string
	RefreshToken string
	TokenType    string
}

////Configurando template
var tpl *template.Template

var banco sync.Map

func init() {
	tpl = template.Must(template.ParseGlob("./view/*.html"))
}

const att = "USER_PASSWORD_AUTH"

///Config template fim

type CognitoConex struct {
	CognitoClient *cognito.CognitoIdentityProvider
	RegFlow       *regFlow
	UserPoolID    string
	AppClientID   string
}

type regFlow struct {
	Username string
}

///Configurações
var conf = &aws.Config{
	Region:                        aws.String("us-east-2"),
	Credentials:                   credentials.AnonymousCredentials,
	MaxRetries:                    aws.Int(1),
	CredentialsChainVerboseErrors: aws.Bool(true),
	//HTTPClient:                    &http.Client{Timeout: 30 * time.Second},
}
var sess, _ = session.NewSession(conf)

///Configurações fim

func main() {

	//creed := credentials.NewStaticCredentials()

	COGNITO_APP_CLIENT_ID := "2va2f2bh1o36q4d2dt6f6qf7k9"
	COGNITO_USER_POOL_ID := "us-east-2_p65qLgx92"

	z := &CognitoConex{
		CognitoClient: cognito.New(sess),
		RegFlow:       &regFlow{},
		UserPoolID:    COGNITO_USER_POOL_ID,
		AppClientID:   COGNITO_APP_CLIENT_ID,
	}

	http.HandleFunc("/registrar", z.Registrar)
	http.HandleFunc("/", home)
	http.HandleFunc("/confirma", z.ConfirmEmail)
	http.HandleFunc("/login", login)
	http.HandleFunc("/auth", z.Aut)
	http.HandleFunc("/test", z.test)
	http.HandleFunc("/deletuser", z.del)
	http.HandleFunc("/admin", z.adm)
	//http.HandleFunc("/test", del)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func (c CognitoConex) adm(w http.ResponseWriter, r *http.Request) {
	c.admAuth()
	fmt.Fprintln(w, "Load...")
}

func (c CognitoConex) del(w http.ResponseWriter, r *http.Request) {
	token := "eyJraWQiOiIxeFE1VFBMaU1tTldBU3hxeGhYdHd4WmJXcnh2eUNiMkxxZmduakpZblh3PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1MjJhMTkyYy1hOWVlLTRkNDUtYjY1OC1iZjg2OWRiMTgyNWMiLCJldmVudF9pZCI6ImFlY2VkNzNhLTk2MTAtNDdhZS04MWVkLTI2NTBmZDMzMzE2OSIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE1ODAzODA3MzcsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTIuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0yX3A2NXFMZ3g5MiIsImV4cCI6MTU4MDM4NDMzNywiaWF0IjoxNTgwMzgwNzM3LCJqdGkiOiIxZTY1YmE1Zi1jYTcwLTRiOTYtODRiNS02YjYxM2VmMzEzY2IiLCJjbGllbnRfaWQiOiIydmEyZjJiaDFvMzZxNGQyZHQ2ZjZxZjdrOSIsInVzZXJuYW1lIjoibGVvIn0.WScBuSE2_lC952XdCCh7Z_1dTfjBTFprGuymr1msJHxVYe4XgTCeIYd4V3DHWtlT8co_5Iqq-TBnzEaN2hR_H9X8vr4rKlKn8essrsHLhFlR-eLhHk6YQ93cq7GQRXPyTMpU-9RTzMGAzhtYbqTqxKY00wyijdA2JFJBX0QVPqhR6S2TrmtRJHJpAUWAfWCV-dlCHf8xmQxAgmWMeVzxIRFYR322QMrv-FwvVmd9HQqIxxa166LEtGhNVKBckwATUwtBBTiBUsycA5rMEnDVgYcLGq2zrxmXRNKAZ8HniFiwQbNtm9JYa0HLtrUw9fPyK88MW01qvfAWQkAjVUCJRA"
	c.deletUser(token)
	fmt.Fprintln(w, "Feiton...")
}

func (c CognitoConex) test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Feito...")

	token := "eyJraWQiOiIxeFE1VFBMaU1tTldBU3hxeGhYdHd4WmJXcnh2eUNiMkxxZmduakpZblh3PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI0ODJiOWQ0Mi0xZTFhLTRiYzMtODcwZC1mN2M5ZDRiMTE2YzIiLCJldmVudF9pZCI6Ijk0OTc5NTIzLTFjNDYtNGU1Yi1hNDYwLTg0YWIxZjZlZGM1MyIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE1ODAzODI3NTIsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTIuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0yX3A2NXFMZ3g5MiIsImV4cCI6MTU4MDM4NjM1MiwiaWF0IjoxNTgwMzgyNzUyLCJqdGkiOiIwMWIyNDQzNi03Y2MwLTRiNDYtYTY1NC01NDkxMTEwZTIyNGUiLCJjbGllbnRfaWQiOiIydmEyZjJiaDFvMzZxNGQyZHQ2ZjZxZjdrOSIsInVzZXJuYW1lIjoibGVvMiJ9.NbcOitAjEE6enGfm4Yj09rKNTmM0YMKjPWvP1LkzWW1irWftcl52Fr2KyaF6qJR6bmw1KLwMZt65R7lXKCuhWO8nh4RPFoc0pu8GsYMe6036gnBrsMh3DVgI6LzPBC-Gk-YOjLN3DF6YWS-BGIxNUvbJebKBhhcYzxVWt-mG0ew2UV083FyI4K1QtZoTbKFgnRZbbb6_CpjpPl024jNoq8K6PqzUtgApp7nXMI2KvZf_NPlW3Q-Q9euxKn-c1EdOlZuIRbToHVRuEz3ATSPq8UwoJxvW9vrLaOC0_5YSApae9lQ_Jx_OpWoJgemIUSPy1SYwmjgRG6Dmnb3LXztunQ"
	c.verificAuth(token)
	//fmt.Fprintln(w, "Retorno: "+vt)
}

func home(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "cadastro.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}

///Tratamento JSON Auth Para OBJETO
func tratamentoAuth(source string) infoAuth {
	var buffer infoAuth

	vt := strings.Split(source, "ExpiresIn: ")
	buffer.ExpiresIn = strings.Split(vt[1], ",")[0]
	vd := strings.Split(source, "\"")
	buffer.AccessToken = vd[1]
	buffer.IdToken = vd[3]
	buffer.RefreshToken = vd[5]
	buffer.TokenType = vd[7]

	return buffer
}
