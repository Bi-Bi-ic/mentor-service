package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rgrs-x/service/api/models/code"
	u "github.com/rgrs-x/service/api/utils"

	"github.com/rgrs-x/service/api/factory"
	"github.com/rgrs-x/service/api/models"
	repository "github.com/rgrs-x/service/api/repository"
	"github.com/rgrs-x/service/api/repository/partner"
	"github.com/rgrs-x/service/api/repository/post"
)

// PostFeatureListController ...
func PostFeatureListController(ctx *gin.Context) {
	var postRepo = post.NewPostRepository(models.GetDB())
	var partnerRepo = partner.NewPartnerRepository(models.GetDB())
	var postFactory factory.PostInfoFactoty
	var partnerEntities = make([]models.Partner, 0)

	postsEntities, getPoststatus := postRepo.GetFeaturePosts()

	if !getPoststatus.AsStatus() {

		response := u.BTResponse{Status: false, Message: getPoststatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	for _, value := range postsEntities {
		partner, getPartnerInfoStatus := partnerRepo.GetDataByID(value.CreatorID)
		if !getPartnerInfoStatus.AsStatus() {
			response := u.BTResponse{Status: false, Message: getPartnerInfoStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		partnerEntities = append(partnerEntities, partner)
	}

	listPostAbles := postFactory.CreateFromListWithPartners(postsEntities, partnerEntities)
	response := u.BTResponse{Status: true, Message: getPoststatus.AsString(), Data: listPostAbles, Code: code.Ok}
	ctx.JSON(http.StatusOK, response)
	return
}

// ListPosts ...
func ListPosts(ctx *gin.Context) {
	var postRepo repository.PostRepository
	var partnerRepo repository.PartnerRepository
	var postFactory factory.PostInfoFactoty

	postRepo = post.NewPostRepository(models.GetDB())
	partnerRepo = partner.NewPartnerRepository(models.GetDB())

	postsEntities, getPoststatus := postRepo.GetPostsList()
	partnerEntities := make([]models.Partner, 0)

	for _, value := range postsEntities {
		partner, getPartnerInfoStatus := partnerRepo.GetDataByID(value.CreatorID)
		if !getPartnerInfoStatus.AsStatus() {
			response := u.BTResponse{Status: false, Message: getPartnerInfoStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		partnerEntities = append(partnerEntities, partner)
	}

	postsListable := postFactory.CreateFromListWithPartners(postsEntities, partnerEntities)

	if getPoststatus.AsStatus() {

		response := u.BTResponse{Status: true, Message: getPoststatus.AsString(), Data: postsListable, Code: code.Ok}
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := u.BTResponse{Status: false, Message: getPoststatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
	ctx.JSON(http.StatusBadRequest, response)
	return

}
