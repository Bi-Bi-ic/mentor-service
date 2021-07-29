package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/comment"
	"github.com/rgrs-x/service/api/repository/course"
	"github.com/rgrs-x/service/api/repository/user"
	u "github.com/rgrs-x/service/api/utils"
)

// CreateComment ...
func CreateComment(ctx *gin.Context) {
	// var postRepo repository.PostRepository
	var userRepo repository.UserRepository

	var content models.Comment

	userID := ctx.Writer.Header().Get("user")
	contentID := ctx.Param("id")
	if userID == "" || contentID == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	bindErr := ctx.ShouldBindJSON(&content)
	if bindErr != nil {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userRepo = user.NewUserRepository(models.GetDB())
	userEntity, getUserByIDStatus := userRepo.GetDataByID(userID)
	if !getUserByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getUserByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	switch content.ContentType {
	case models.PostConent:
		return

	case models.CourseContent:
		createCourseComment(ctx, userEntity, content, contentID, userID)
		return

	default:
		response := u.BTResponse{Status: false, Message: "Content-Type Invalid", Data: []string{}, Code: code.ResourceError}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
}

func createCourseComment(ctx *gin.Context, userEntity models.User, content models.Comment, contentID, userID string) {
	var commentFactory factory.CommentFactory
	var userFactory factory.UserInfoFactory

	var courseRepo = course.NewCourseRepository(models.GetDB())

	customerAdvance, getCourseByIDStatus := courseRepo.GetCourseById(contentID)
	if !getCourseByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCourseByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	commentEntity, createdStatus := courseRepo.CreateComment(userEntity.ID.String(), customerAdvance, content)

	if createdStatus.AsStatus() {

		userable := userFactory.Create(userEntity)
		commentable := commentFactory.Create(commentEntity, userable)

		response := u.BTResponse{Status: true, Message: createdStatus.AsString(), Data: commentable, Code: code.Created}

		ctx.JSON(http.StatusOK, response)
		return

	}

	response := u.BTResponse{Status: false, Message: createdStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
	ctx.JSON(http.StatusBadRequest, response)
	return
}

// GetCourseComments ...
func GetCourseComments(ctx *gin.Context) {

	var courseRepo = course.NewCourseRepository(models.GetDB())
	var userRepo = user.NewUserRepository(models.GetDB())
	var commentFactory factory.CommentFactory
	var userFactory factory.UserInfoFactory
	var commentsAble = []factory.CommentAble{}

	contentID := ctx.Param("id")
	if contentID == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	commentEntities, getCommentStatus := courseRepo.GetComment(contentID)
	if !getCommentStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCommentStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	for _, content := range commentEntities {
		userEntity, getUserByIDStatus := userRepo.GetDataByID(content.UserID)
		if !getUserByIDStatus.AsStatus() {
			response := u.BTResponse{Status: false, Message: getUserByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		userAble := userFactory.Create(userEntity)
		customerAdvance := commentFactory.Create(content, userAble)

		commentsAble = append(commentsAble, customerAdvance)
	}

	response := u.BTResponse{Status: true, Message: getCommentStatus.AsString(), Data: commentsAble, Code: code.Ok}
	ctx.JSON(http.StatusOK, response)
	return

}

// UpdateComment ...
func UpdateComment(ctx *gin.Context) {
	var userRepo = user.NewUserRepository(models.GetDB())
	var commentRepo = comment.NewCommentRepository(models.GetDB())

	var commentFactory factory.CommentFactory
	var userFactory factory.UserInfoFactory

	var payload models.Comment

	creatorID := ctx.Writer.Header().Get("user")
	commentID := ctx.Param("id")
	if creatorID == "" || commentID == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	commentEntity, getCommentByIDStatus := commentRepo.GetCommentByID(commentID)
	if !getCommentByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCommentByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if commentEntity.UserID != creatorID {
		response := u.BTResponse{Status: false, Message: repository.CanNotUpdate.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userEntity, getUserByIDStatus := userRepo.GetDataByID(creatorID)
	if !getUserByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getUserByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	bindErr := ctx.ShouldBindJSON(&payload)
	if bindErr != nil {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if payload.ContentType != "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	commentEntity.Content = payload.Content
	customerAdvance, updateCommentStatus := commentRepo.UpdateComment(commentEntity)
	if updateCommentStatus.AsStatus() {
		userAble := userFactory.Create(userEntity)
		commentAble := commentFactory.Create(customerAdvance, userAble)

		response := u.BTResponse{Status: true, Message: updateCommentStatus.AsString(), Data: commentAble, Code: code.Ok}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := u.BTResponse{Status: false, Message: updateCommentStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
	ctx.JSON(http.StatusBadRequest, response)
	return
}

// DeleteComment ...
func DeleteComment(ctx *gin.Context) {
	var userRepo = user.NewUserRepository(models.GetDB())
	var commentRepo = comment.NewCommentRepository(models.GetDB())

	creatorID := ctx.Writer.Header().Get("user")
	commentID := ctx.Param("id")
	if creatorID == "" || commentID == "" {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	commentEntity, getCommentByIDStatus := commentRepo.GetCommentByID(commentID)
	if !getCommentByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCommentByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	if commentEntity.UserID != creatorID {
		response := u.BTResponse{Status: false, Message: repository.CanNotDelete.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, getUserByIDStatus := userRepo.GetDataByID(creatorID)
	if !getUserByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getUserByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	deleteCommentStatus := commentRepo.DeleteComment(commentEntity)
	if deleteCommentStatus.AsStatus() {
		response := u.BTResponse{Status: true, Message: message.DeletedResource, Data: []string{}, Code: code.Deleted}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := u.BTResponse{Status: false, Message: deleteCommentStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
	ctx.JSON(http.StatusBadRequest, response)
	return
}
