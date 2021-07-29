package user

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/rgrs-x/service/api/models/code"
	u "github.com/rgrs-x/service/api/utils"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/user"
)

func UserListController(c *gin.Context) {
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory

	userRepo = user.NewUserRepository(models.GetDB())

	userEntities, getUsersStatus := userRepo.GetAllUsers()

	userListable := userFactory.CreateFromList(userEntities)

	if getUsersStatus.AsStatus() {

		response := u.BTResponse{Status: true, Message: getUsersStatus.AsString(), Data: userListable, Code: code.Ok}

		c.JSON(http.StatusOK, response)
	} else {
		response := u.BTResponse{Status: false, Message: getUsersStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
	}

}

// UserFeatureListController ...
func UserFeatureListController(c *gin.Context) {
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory

	userRepo = user.NewUserRepository(models.GetDB())

	userEntities, getUsersStatus := userRepo.GetFeatureUsers()

	userListable := make([]factory.Userable, 0)
	userListable = userFactory.CreateFromList(userEntities)

	if getUsersStatus.AsStatus() {

		response := u.BTResponse{Status: true, Message: getUsersStatus.AsString(), Data: userListable, Code: code.Ok}

		c.JSON(http.StatusOK, response)
	} else {
		response := u.BTResponse{Status: false, Message: getUsersStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
	}

}
