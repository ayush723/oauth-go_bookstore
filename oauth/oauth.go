package oauth

import (
	"net/http"

	"github.com/ayush723/oauth-go_bookstore/oauth/errors"
)


const(
	headerXPublic = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-User-Id"
)

type oauthClient struct{}

type oauthInterface interface{}

func IsPublic(request *http.Request) bool{
	if request == nil{
		return true //means it is a public request
	}
	return request.Header.Get(headerXPublic) == "true"

}

func AuthenticateRequest(request *http.Request) *errors.RestErr{
	if request == nil{
		return  //means it is a public request
	}
}