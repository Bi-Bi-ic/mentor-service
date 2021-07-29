package course

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/rgrs-x/service/api/models/code"
	u "github.com/rgrs-x/service/api/utils"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/course"
	"github.com/rgrs-x/service/api/repository/user"
)

func CourseListController(c *gin.Context) {
	var courseFactory factory.CourseInfoFactory
	var courseRepo repository.CourseRepository
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory

	courseFactory = factory.CourseInfoFactory{}
	userFactory = factory.UserInfoFactory{}

	courseRepo = course.NewCourseRepository(models.GetDB())
	userRepo = user.NewUserRepository(models.GetDB())

	courseListable := make([]factory.Courseable, 0)
	// userEntityList := make([]models.User, 0)
	customerAdvance, getAllStatus := courseRepo.GetAllCourses()

	for _, value := range customerAdvance {
		userEntity, _ := userRepo.GetUserByID(value.UserID)
		userable := userFactory.Create(userEntity)
		courseable := courseFactory.Create(value, userable)

		courseListable = append(courseListable, courseable)
	}

	if getAllStatus.AsStatus() {

		response := u.BTResponse{Status: true, Message: getAllStatus.AsString(), Data: courseListable, Code: code.Ok}

		c.JSON(http.StatusOK, response)
	} else {
		response := u.BTResponse{Status: false, Message: getAllStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
	}

}

func CourseFeatureListController(c *gin.Context) {
	var courseRepo repository.CourseRepository

	courseRepo = course.NewCourseRepository(models.GetDB())

	customerAdvance, getAllStatus := courseRepo.GetCoursesFeature()

	if getAllStatus.AsStatus() {
		courseListable := make([]factory.Courseable, 0)
		courseListable = GetCourseListable(customerAdvance)
		response := u.BTResponse{Status: true, Message: getAllStatus.AsString(), Data: courseListable, Code: code.Ok}

		c.JSON(http.StatusOK, response)
	} else {
		response := u.BTResponse{Status: false, Message: getAllStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
	}
}

// GetCourseListable ...
func GetCourseListable(courseEnities []models.CourseEntity) []factory.Courseable {
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory
	var courseFactory factory.CourseInfoFactory
	userRepo = user.NewUserRepository(models.GetDB())

	courseListable := make([]factory.Courseable, 0)
	userFactory = factory.UserInfoFactory{}
	courseFactory = factory.CourseInfoFactory{}

	for _, value := range courseEnities {
		userEntity, _ := userRepo.GetUserByID(value.UserID)
		userable := userFactory.Create(userEntity)
		courseable := courseFactory.Create(value, userable)

		courseListable = append(courseListable, courseable)
	}

	return courseListable
}
