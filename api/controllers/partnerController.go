package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/rgrs-x/service/api/models"
	"github.com/rgrs-x/service/api/models/code"
	"github.com/rgrs-x/service/api/repository/partner"
	"github.com/rgrs-x/service/api/factory"
	u "github.com/rgrs-x/service/api/utils"
)

// LikeMentor ...
func LikeMentor(c *gin.Context) {
	mentorID := c.Param("id")

	var partnerRepository = partner.NewPartnerRepository(models.GetDB())
	repo, status := partnerRepository.UpdateMentorLike(mentorID)

	handlerStatus(repo, status, models.PartnerNormal, c)
}
// PublicPartnerInfo ...
func PublicPartnerInfo(c *gin.Context) {
	var partnerFactory factory.PartnerInfoFactory

	partnerFactory = factory.PartnerInfoFactory{}

	partnerID := c.Param("id")
	
	partnerRepository := partner.NewPartnerRepository(models.GetDB())
	partnerEntity, getPartnerInfoStatus := partnerRepository.GetPartnerInfo(partnerID)

	partnerPublish := partnerFactory.CreateFromModel(partnerEntity)

	if getPartnerInfoStatus.AsStatus() {
		
		response := u.BTResponse{Status: true, Message: getPartnerInfoStatus.AsString(), Data: partnerPublish, Code: code.Ok}

		c.JSON(http.StatusOK, response)	
	} else {
		response := u.BTResponse{Status: false, Message: getPartnerInfoStatus.AsString(), Data: []string{}, Code: code.DataIsEmpty}
		c.JSON(http.StatusBadRequest, response)
	}

}
