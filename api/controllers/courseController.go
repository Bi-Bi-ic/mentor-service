package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/course"
	"github.com/rgrs-x/service/api/repository/user"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"
)

// CreateCourse ...
func CreateCourse(c *gin.Context) {

	var customer models.CourseEntity
	var courseRepo repository.CourseRepository
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory
	var courseFactory factory.CourseInfoFactory

	courseRepo = course.NewCourseRepository(models.GetDB())
	userRepo = user.NewUserRepository(models.GetDB())
	userFactory = factory.UserInfoFactory{}
	courseFactory = factory.CourseInfoFactory{}
	customer = models.CourseEntity{}

	err := c.ShouldBindJSON(&customer)
	userID := c.Writer.Header().Get("user")

	customer.UserID = userID

	if err != nil {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	customerAdvance, _ := courseRepo.CreateCourse(customer)

	userEntity, _ := userRepo.GetUserByID(userID)

	userable := userFactory.Create(userEntity)

	courseable := courseFactory.Create(customerAdvance, userable)

	c.JSON(http.StatusOK, courseable)

}

// UpdateCourse ...
func UpdateCourse(c *gin.Context) {

	var customer models.CourseEntity
	var courseRepo repository.CourseRepository
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory
	var courseFactory factory.CourseInfoFactory

	courseRepo = course.NewCourseRepository(models.GetDB())
	userRepo = user.NewUserRepository(models.GetDB())

	courseFactory = factory.CourseInfoFactory{}
	customer = models.CourseEntity{}

	customer = models.CourseEntity{}

	courseID, err := uuid.FromString(c.Param("id"))

	if err != nil {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	bindErr := c.ShouldBindJSON(&customer)
	customer.ID = courseID

	userID := c.Writer.Header().Get("user")
	customer.UserID = userID

	if bindErr != nil {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	customerAdvance, _ := courseRepo.UpdateCourse(customer)
	userEntity, _ := userRepo.GetUserByID(userID)

	userable := userFactory.Create(userEntity)

	courseable := courseFactory.Create(customerAdvance, userable)

	c.JSON(http.StatusOK, courseable)

}

// GetAllCourseByMentorId ...
func GetAllCourseByMentorId(c *gin.Context) {
	var courseFactory factory.CourseInfoFactory
	var courseRepo repository.CourseRepository
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory

	courseFactory = factory.CourseInfoFactory{}
	userFactory = factory.UserInfoFactory{}

	courseRepo = course.NewCourseRepository(models.GetDB())
	userRepo = user.NewUserRepository(models.GetDB())

	courseable := make([]factory.Courseable, 0)

	userId := c.Param("id")

	if userId == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	customerAdvance, _ := courseRepo.GetAllCourseByMentorId(userId)
	userEntity, _ := userRepo.GetUserByID(userId)

	userable := userFactory.Create(userEntity)
	factoryCourseableOutput := courseFactory.CreateFromList(customerAdvance, userable)

	if factoryCourseableOutput != nil {
		courseable = factoryCourseableOutput
	}

	c.JSON(http.StatusOK, courseable)
}

// GetAllCourseByMentorId ...
func GetCourseById(c *gin.Context) {
	var courseFactory factory.CourseInfoFactory
	var courseRepo repository.CourseRepository
	var userRepo repository.UserRepository
	var userFactory factory.UserInfoFactory

	courseFactory = factory.CourseInfoFactory{}
	userFactory = factory.UserInfoFactory{}

	courseRepo = course.NewCourseRepository(models.GetDB())
	userRepo = user.NewUserRepository(models.GetDB())

	courseId := c.Param("id")

	if courseId == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	customerAdvance, getCourseByIdStatus := courseRepo.GetCourseById(courseId)
	userEntity, _ := userRepo.GetUserByID(customerAdvance.UserID)

	userable := userFactory.Create(userEntity)
	factoryCourseableOutput := courseFactory.Create(customerAdvance, userable)

	if getCourseByIdStatus.AsStatus() {

		response := u.BTResponse{Status: true, Message: getCourseByIdStatus.AsString(), Data: factoryCourseableOutput, Code: code.Ok}

		c.JSON(http.StatusOK, response)
	} else {
		response := u.BTResponse{Status: false, Message: getCourseByIdStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
	}
}

// GetAllCourseByMentorId ...
func DeleteCourseById(c *gin.Context) {
	var courseRepo repository.CourseRepository
	userId := c.Writer.Header().Get("user")
	courseId := c.Param("id")

	courseRepo = course.NewCourseRepository(models.GetDB())

	_, deleteStatus := courseRepo.DeleteCourseById(courseId, userId)

	if deleteStatus.AsStatus() {
		c.JSON(http.StatusOK, u.Message(true, deleteStatus.AsString()))
	} else {
		c.JSON(http.StatusBadRequest, u.Message(false, deleteStatus.AsString()))
	}

}

// RegisterCourse ...
func RegisterCourse(ctx *gin.Context) {
	var courseRepo = course.NewCourseRepository(models.GetDB())
	var userRepo = user.NewUserRepository(models.GetDB())

	courseID := ctx.Param("id")
	userID := ctx.Writer.Header().Get("user")
	if userID == "" || courseID == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, getUserByIDStatus := userRepo.GetDataByID(userID)
	if !getUserByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getUserByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	courseEntity, getCourseByIDStatus := courseRepo.GetCourseById(courseID)
	if !getCourseByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCourseByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	registerCourseStatus := courseRepo.RegisterCourse(courseEntity, userID)
	if !registerCourseStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: registerCourseStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := u.BTResponse{Status: true, Message: registerCourseStatus.AsString(), Data: []string{}, Code: code.Ok}
	ctx.JSON(http.StatusOK, response)
	return
}

// GetMenteesFromCourse ...
func GetMenteesFromCourse(ctx *gin.Context) {
	var courseRepo = course.NewCourseRepository(models.GetDB())
	var userRepo = user.NewUserRepository(models.GetDB())

	var menteeEntities []models.User
	var menteeFactory factory.MenteeInfoFactory

	courseID := ctx.Param("id")
	listMenteesID, getMenteesStatus := courseRepo.GetMentees(courseID)
	if !getMenteesStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getMenteesStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	for _, value := range listMenteesID {
		userEntity, getDataByIDStatus := userRepo.GetDataByID(value)
		if !getDataByIDStatus.AsStatus() {
			response := u.BTResponse{Status: false, Message: getDataByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		menteeEntities = append(menteeEntities, userEntity)
	}

	menteeAble := menteeFactory.CreateFromList(menteeEntities)
	response := u.BTResponse{Status: true, Message: getMenteesStatus.AsString(), Data: menteeAble, Code: code.Ok}
	ctx.JSON(http.StatusOK, response)
	return

}

// LikeCourse ...
func LikeCourse(ctx *gin.Context) {

	var courseRepo repository.CourseRepository

	courseRepo = course.NewCourseRepository(models.GetDB())

	courseID := ctx.Param("id")

	customerAdvance, getCourseByIdStatus := courseRepo.GetCourseById(courseID)

	userID := ctx.Writer.Header().Get("user")

	if userID == customerAdvance.UserID {
		response := u.BTResponse{Status: false, Message: "Your creator this course, and you can not like yourself", Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	customerAdvance.TotalLike = customerAdvance.TotalLike + 1

	if !getCourseByIdStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCourseByIdStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
	} else {
		_, updateLikeStatus := courseRepo.UpdateCourse(customerAdvance)

		if updateLikeStatus.AsStatus() {
			response := u.BTResponse{Status: true, Message: updateLikeStatus.AsString(), Data: []string{}, Code: code.Ok}
			ctx.JSON(http.StatusOK, response)
		}
	}

}
