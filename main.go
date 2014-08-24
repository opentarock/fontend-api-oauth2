package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/arjantop/oauth2-util"
	"github.com/opentarock/service-api/go/client"
	"github.com/opentarock/service-api/go/proto_oauth2"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	r := gin.Default()

	oauth2Service, err := client.NewOauth2ClientNanomsg()
	if err != nil {
		log.Fatalf("Error creating oauth2 client: %s", err)
	}

	r.POST("/token", func(c *gin.Context) {
		paramGrantType := c.Request.PostFormValue(oauth2.ParameterGrantType)
		paramUsername := c.Request.PostFormValue(oauth2.ParameterUsername)
		paramPassword := c.Request.PostFormValue(oauth2.ParameterPassword)
		paramRefreshToken := c.Request.PostFormValue(oauth2.ParameterRefreshToken)

		client, err := readBasicAuth(c.Request)
		if err != nil {
			c.Writer.Header().Set("WWW-Authenticate", `Basic realm="oauth2"`)
			c.Abort(http.StatusUnauthorized)
			return
		}

		request := &proto_oauth2.AccessTokenRequest{
			GrantType:    checkEmpty(paramGrantType),
			Username:     checkEmpty(paramUsername),
			Password:     checkEmpty(paramPassword),
			RefreshToken: checkEmpty(paramRefreshToken),
		}

		response, err := oauth2Service.GetAccessToken(client.GetId(), client.GetSecret(), request)
		if err != nil {
			c.Abort(http.StatusInternalServerError)
			return
		}

		if response.GetSuccess() {
			c.JSON(http.StatusOK, response.GetToken())
		} else {
			oauthErr := response.GetError()
			httpCode := oauth2.ErrorHttpResponseCode(oauthErr.GetError())
			c.JSON(httpCode, oauthErr)
		}
	})
	r.Run(":8080")
}

// readBasicAuth reads an Authorization header using Basic access authentication method.
// On Success client with its id and secret read is returned.
// Error is returned if the header is missing, wrong authentication method is used or its
// format is invalid.
func readBasicAuth(r *http.Request) (*proto_oauth2.Client, error) {
	auth := r.Header.Get("Authorization")
	authParts := strings.SplitN(auth, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Basic" {
		return nil, errors.New("Authorization method must be Basic")
	}
	decoded, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		return nil, err
	}
	clientCredentialsParts := strings.SplitN(string(decoded), ":", 2)
	if len(clientCredentialsParts) != 2 {
		return nil, errors.New("Client credentials format error")
	}
	return &proto_oauth2.Client{Id: checkEmpty(clientCredentialsParts[0]), Secret: checkEmpty(clientCredentialsParts[1])}, nil
}

// checkEmpty converts a string value to a string pointer.
// Nil is returned if value is empty, otherwise pointer
// to the string value is returned.
func checkEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
