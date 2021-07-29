package repository

import (
	"github.com/rgrs-x/service/api/models"
)

//PartnerRepository is an interface can be implemented
type PartnerRepository interface {
	GetDataByID(string) (models.Partner, Status)
	UpdateMentorLike(string) (RepoResponse, Status)
	GetPartnerInfo(id string) (models.Partner, Status)

	// Using for Contents
	CountPosts(*models.Partner) error
	CountPostsExpired(*models.Partner) error
}
