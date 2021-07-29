package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	courseControl "github.com/rgrs-x/service/api/controllers/course"
	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/search"
	u "github.com/rgrs-x/service/api/utils"
)

func SearchUser(c *gin.Context) {

	var searchRepo repository.SearchRepositoryable
	var userFactory factory.UserInfoFactory

	searchRepo = search.NewSearchRepository(models.GetDB())
	userFactory = factory.UserInfoFactory{}

	userName := c.Query("q")

	if userName == "" {
		response := u.Message(true, "Catnot find something with empty string")
		response["data"] = make([]factory.Userable, 0)
		c.JSON(http.StatusOK, response)
		return
	}

	users, seachStatus := searchRepo.SearchUserName(userName)

	userable := userFactory.CreateFromList(users)

	if seachStatus.AsStatus() {
		response := u.Message(true, seachStatus.AsString())
		response["data"] = userable
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, u.Message(false, seachStatus.AsString()))
	}

}

// SearchPost ...
func SearchPost(c *gin.Context) {

	var searchRepo repository.SearchRepositoryable
	var postFactory factory.PostInfoFactoty

	searchRepo = search.NewSearchRepository(models.GetDB())
	postFactory = factory.PostInfoFactoty{}

	textReSearch := c.Query("q")

	if textReSearch == "" {
		response := u.Message(true, "Catnot find something with empty string")
		response["data"] = make([]factory.Postable, 0)
		c.JSON(http.StatusOK, response)
		return
	}

	postsListable := make([]factory.Postable, 0)

	posts, seachStatus := searchRepo.SearchPost(textReSearch)

	postsListable = postFactory.CreateFromList(posts)

	if seachStatus.AsStatus() {
		response := u.Message(true, seachStatus.AsString())
		response["data"] = postsListable
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, u.Message(false, seachStatus.AsString()))
	}

}
// SearchCourse ...
func SearchCourse(c *gin.Context) {
	var searchRepo repository.SearchRepositoryable
	courseListable := make([]factory.Courseable, 0)

	searchRepo = search.NewSearchRepository(models.GetDB())

	textReSearch := c.Query("q")

	if textReSearch == "" {
		response := u.Message(true, "Catnot find something with empty string")
		response["data"] = make([]factory.Courseable, 0)
		c.JSON(http.StatusOK, response)
		return
	}

	courses, seachStatus := searchRepo.SearchCourse(textReSearch)
	courseListable = courseControl.GetCourseListable(courses)

	if seachStatus.AsStatus() {
		response := u.Message(true, seachStatus.AsString())
		response["data"] = courseListable
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, u.Message(false, seachStatus.AsString()))
	}

}
// SearchAll ...
func SearchAll(c *gin.Context) {
	var searchRepo repository.SearchRepositoryable
	var userFactory factory.UserInfoFactory
	var postFactory factory.PostInfoFactoty

	searchRepo = search.NewSearchRepository(models.GetDB())
	userFactory = factory.UserInfoFactory{}
	postFactory = factory.PostInfoFactoty{}

	textReSearch := c.Query("q")

	users, seachStatus := searchRepo.SearchUserName(textReSearch)
	posts, _ := searchRepo.SearchPost(textReSearch)
	courses, _ := searchRepo.SearchCourse(textReSearch)

	courseListable := make([]factory.Courseable, 0)
	postsListable := make([]factory.Postable, 0)
	userable := make([]factory.Userable, 0)

	userable = userFactory.CreateFromList(users)
	courseListable = courseControl.GetCourseListable(courses)
	postsListable = postFactory.CreateFromList(posts)

	if seachStatus.AsStatus() {
		dataAdvance := map[string]interface{}{}
		dataAdvance["users"] = userable
		dataAdvance["posts"] = postsListable
		dataAdvance["courses"] = courseListable
		response := u.Message(true, seachStatus.AsString())
		response["data"] = dataAdvance
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, u.Message(false, seachStatus.AsString()))
	}

}
