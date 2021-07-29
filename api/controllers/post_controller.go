package controllers

import (
	"net/http"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/models/message"
	"github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/partner"
	"github.com/rgrs-x/service/api/repository/post"
	u "github.com/rgrs-x/service/api/utils"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

// PostHandler ...
func PostHandler(statusCode int, Data u.ResultRepository, c *gin.Context) {
	var postFactory factory.PostInfoFactoty
	c.Writer.Header().Set("Content-Type", "application/json")

	switch statusCode {
	case http.StatusOK, http.StatusCreated:
		Data.Result = postFactory.NewPost(Data.Result)
		c.JSON(statusCode, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.Contents:
		Data.Result = postFactory.NewListPost(Data.Result)
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	case models.TrackingSuccess:
		c.JSON(http.StatusOK, u.Response{Status: true, Message: Data.Message, Data: Data.Result})
		return
	default:
		c.JSON(statusCode, u.Response{Status: false, Message: Data.Error.Error(), Data: Data.Result})
		return
	}
}

// CreatePost make a post request
func CreatePost(c *gin.Context) {
	//We init a new Post sruct here
	draft := models.Post{}
	tempID := c.Writer.Header().Get("user")

	//Then Bind the Request to new Post
	err := c.ShouldBindJSON(&draft)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: "Invalid Request", Data: []string{}})
		return
	}

	//And call the Post-Repository to create Post
	postRepository := post.NewPostRepository(models.GetDB())
	resp, statusCode := postRepository.Create(draft, tempID)

	PostHandler(statusCode, resp, c)
}

// GetPost find an existing Post
func GetPost(ctx *gin.Context) {
	var postRepo = post.NewPostRepository(models.GetDB())
	var partnerRepo = partner.NewPartnerRepository(models.GetDB())

	var postFactory factory.PostInfoFactoty

	contentID := ctx.Param("id")
	postEntity, getPostByIDStatus := postRepo.GetPostDetails(contentID)
	if !getPostByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getPostByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	creatorEntity, getCreatorByIDStatus := partnerRepo.GetDataByID(postEntity.CreatorID)
	if !getCreatorByIDStatus.AsStatus() {
		response := u.BTResponse{Status: false, Message: getCreatorByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	postAble := postFactory.CreatedWithCreator(postEntity, creatorEntity)
	response := u.BTResponse{Status: true, Message: getCreatorByIDStatus.AsString(), Data: postAble, Code: code.Ok}
	ctx.JSON(http.StatusOK, response)
	return
}

// UpdatePost change an existing Post with id received
func UpdatePost(c *gin.Context) {
	id := c.Params.ByName("id")
	draft := models.Post{}

	creator := c.Writer.Header().Get("user")
	err := c.ShouldBindJSON(&draft)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request")
		return
	}

	postRepository := post.NewPostRepository(models.GetDB())
	resp, statusCode := postRepository.UpdatePost(draft, id, creator)

	if statusCode == http.StatusOK {
		var factory factory.PostInfoFactoty
		resp["data"] = factory.NewPost(resp["data"])

		c.JSON(statusCode, resp)
	} else {
		c.JSON(statusCode, resp)
	}
}

// DeletePost delete an existing Post with id received
func DeletePost(c *gin.Context) {
	id := c.Params.ByName("id")
	creator := c.Writer.Header().Get("user")

	postRepository := post.NewPostRepository(models.GetDB())
	resp, statusCode := postRepository.DeletePost(id, creator)

	c.JSON(statusCode, resp)
}

// LikePost ...
func LikePost(ctx *gin.Context) {
	var postRepo = post.NewPostRepository(models.GetDB())
	var postFactory factory.PostInfoFactoty

	contentID := ctx.Params.ByName("id")
	contentEntity, getPostByIDStatus := postRepo.GetPostByID(contentID)
	if !getPostByIDStatus.AsStatus() {

		response := u.BTResponse{Status: false, Message: getPostByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	customAdvance, updatePostLikeStatus := postRepo.UpdatePostLike(contentEntity)
	if updatePostLikeStatus.AsStatus() {
		postAble := postFactory.Create(customAdvance)

		response := u.BTResponse{Status: true, Message: updatePostLikeStatus.AsString(), Data: postAble, Code: code.Ok}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := u.BTResponse{Status: false, Message: updatePostLikeStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
	ctx.JSON(http.StatusBadRequest, response)
	return

}

// ReadPost ...
func ReadPost(ctx *gin.Context) {
	var postRepo = post.NewPostRepository(models.GetDB())
	var postFactory factory.PostInfoFactoty

	contentID := ctx.Params.ByName("id")
	contentEntity, getPostByIDStatus := postRepo.GetPostByID(contentID)
	if !getPostByIDStatus.AsStatus() {

		response := u.BTResponse{Status: false, Message: getPostByIDStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	customAdvance, updatePostReviewStatus := postRepo.UpdatePostReview(contentEntity)
	if updatePostReviewStatus.AsStatus() {
		postAble := postFactory.Create(customAdvance)

		response := u.BTResponse{Status: true, Message: updatePostReviewStatus.AsString(), Data: postAble, Code: code.Ok}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := u.BTResponse{Status: false, Message: updatePostReviewStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
	ctx.JSON(http.StatusBadRequest, response)
	return
}

// CreateIntroductionPost ...
func CreateIntroductionPost(c *gin.Context) {
	draft := models.Post{}
	creatorID := c.Writer.Header().Get("user")
	_, err := uuid.FromString(creatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	//Then Bind the Request to new Post
	err = c.ShouldBindJSON(&draft)
	if err != nil {
		response := u.BTResponse{Status: false, Message: message.BadRequest, Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)

		return
	}

	//And call the Post-Repository to create Post
	postRepository := post.NewPostRepository(models.GetDB())
	repo, status := postRepository.CreateIntroduction(draft, creatorID)

	handlerStatus(repo, status, models.IntroductionNormal, c)
}

// GetIntroductionPost ...
func GetIntroductionPost(c *gin.Context) {
	creatorID := c.Writer.Header().Get("user")
	_, err := uuid.FromString(creatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, u.Response{Status: false, Message: repository.ErrMalformedID.Error(), Data: []string{}})
		return
	}

	//And call the Post-Repository to create Post
	postRepository := post.NewPostRepository(models.GetDB())
	repo, status := postRepository.GetIntroduction(creatorID)

	handlerStatus(repo, status, models.IntroductionNormal, c)
}
