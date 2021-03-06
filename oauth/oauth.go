package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ayush723/oauth-go_bookstore/oauth/errors"
	"github.com/mercadolibre/golang-restclient/rest"
)




const(
	headerXPublic = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-Caller-Id"

	paramAccessToken = "access_token"

)

var(
	oauthRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 200 * time.Millisecond,
	}
)

type accessToken struct{
	Id string 	`json:"id"`
	UserId int64 `json:"user_id"`
	ClientId int16 `json:"client_id"`
}

// type oauthInterface interface{}

func IsPublic(request *http.Request) bool{
	if request == nil{
		return true //means it is a public request
	}
	return request.Header.Get(headerXPublic) == "true"

}
//GetCallerId returns caller id from header
func GetCallerId(request *http.Request)int64{
	if request == nil{
		return 0
	}
	callerId, err := strconv.ParseInt(request.Header.Get(headerXCallerId), 10,64)
	if err != nil{
		return 0
	} 
	return callerId

}

func GetClientId(request *http.Request)int64{
	if request == nil{
		return 0
	}
	clientId, err := strconv.ParseInt(request.Header.Get(headerXClientId), 10,64)
	if err != nil{
		return 0
	} 
	return clientId

}

func AuthenticateRequest(request *http.Request) *errors.RestErr{
	if request == nil{
		return  nil//means it is a public request
		}
		//clears headers if present
		cleanRequest(request)

		//gets access token id
		accessTokenId:= strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
		//http://api.bookstore.com/resource?access_token=abc123
		if accessTokenId == ""{
			return nil
		}

		//hits the oauth api/access_token/:access_token_id
		at, err := getAccessToken(accessTokenId)
		if err != nil{
			if err.Status == http.StatusNotFound{
				return nil
			}
			return err
		}
		//adds the following headers
		request.Header.Add(headerXClientId, fmt.Sprintf("%v",at.ClientId))
		request.Header.Add(headerXCallerId, fmt.Sprintf("%v",at.UserId))

		return nil
}

func cleanRequest(request *http.Request){
	if request == nil{
		return
	}
	request.Header.Del(headerXClientId)
	request.Header.Del(headerXCallerId)
}

func getAccessToken(accessTokenId string) (*accessToken, *errors.RestErr){
	response := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenId))
	if response == nil || response.Request == nil {
		return nil, errors.NewInternalServerError("invalid restclient response when trying to get access token")

	}
	if response.StatusCode > 299{
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil{
			return nil, errors.NewInternalServerError("invalid error interface when trying to get access token")
		}
		
		return nil, &restErr
	}
	// if response is ok, we unmarshal the response body
	var at accessToken
	if err := json.Unmarshal(response.Bytes(), &at); err != nil{
		return nil, errors.NewInternalServerError("error when trying to unmarshal access token response")
	}
	return &at, nil
}