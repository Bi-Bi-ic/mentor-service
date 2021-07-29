package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/repository/user"
	u "github.com/rgrs-x/service/api/utils"
)

// RefreshToken ...
func RefreshToken(c *gin.Context) *models.Token {

	//@ Grab the token from the header
	tokenHeader := c.Request.Header.Get("Authorization")

	//@ Token is missing, returns with error code 403 Unauthorized
	if tokenHeader == "" {
		return nil
	}

	//@ The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return nil
	}

	//@ Grab the token part, what we are truly interested in
	tokenPart := splitted[1]
	tk := &models.Token{}

	_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	//Well... Not the token wanted --> http code 403 bye bye
	if tk.Type != models.Refresh {
		return nil
	}

	//@ Malformed token with error signature, returns with http code 403 as usual
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil
		}
	}

	return tk
}

// UserToken ...
func UserToken(c *gin.Context) {
	resp := RefreshToken(c)
	if resp == nil {
		result := u.Message(false, "Invalid or Malformed auth token")
		c.JSON(http.StatusForbidden, result)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Abort()
		return
	}

	customer := models.User{}
	userRepository := user.NewUserRepository(models.GetDB())

	response, statusCode := userRepository.GetByID(resp.UserId.String(), customer)
	if statusCode == http.StatusOK {
		// Reflect with user entity
		user := models.User{}
		err := mapstructure.Decode(response["data"], &user)
		if err != nil {
			panic(err)
		}

		user.GenerateToken()
		//@ delete password
		user.Password = ""

		result := u.Message(true, "Token has been refreshed")

		var userFactory factory.UserInfoFactory
		result["data"] = userFactory.Create(user)

		c.JSON(http.StatusOK, result)
		c.Abort()
		c.Writer.Header().Set("Content-Type", "application/json")
		return
	}

	response = u.Message(false, "Token is not valid")
	c.JSON(http.StatusForbidden, response)
	return
}
