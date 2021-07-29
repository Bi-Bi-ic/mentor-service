package partner

import (
	"errors"
	"strconv"
	"time"

	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	"gorm.io/gorm"
)

// partnerStorage is a struct implementing PartnerRepository interface{}
type partnerStorage struct {
	Db *gorm.DB
}

// NewPartnerRepository ... We can implement to use PartnerRepository interface{} there
func NewPartnerRepository(db *gorm.DB) repo.PartnerRepository {
	return &partnerStorage{
		Db: db,
	}
}

func (storage *partnerStorage) GetPartnerInfo(id string) (partner models.Partner, status repo.Status) {
	mentor, err := storage.getPartnerProfile(id)
	if err != nil {
		status = repo.NotFound
		return
	}
	partner = mentor
	status = repo.Success
	return
}

// GetDataByID ...
func (storage *partnerStorage) GetDataByID(partnerID string) (partner models.Partner, status repo.Status) {
	queryStmt := storage.Db.Where("id = ?", partnerID)
	find := queryStmt.First(&partner)
	err := find.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = repo.NotFound
			return
		}

		status = repo.CannotGet
		return
	}

	err = storage.CountPosts(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostAvailable = "0"
		} else {
			status = repo.CannotGet
			return
		}
	}

	err = storage.CountPostsExpired(&partner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			partner.PostExpired = "0"
		} else {
			status = repo.CannotGet
			return
		}
	}

	status = repo.Success
	return
}

// Get Partner Profile ...
func (storage *partnerStorage) getPartnerProfile(partnerID string) (partner models.Partner, err error) {
	queryStmt := storage.Db.Table("partners").Where("id = ?", partnerID).First(&partner)
	err = queryStmt.Error
	if err != nil {
		return
	}

	return
}

// CountPosts ...
func (storage *partnerStorage) CountPosts(partner *models.Partner) error {
	var count int64

	queryStatement := storage.Db.Model(&models.Post{}).
		Where("creator_id = ?", partner.ID).
		Count(&count)

	err := queryStatement.Error
	partner.PostAvailable = strconv.FormatInt(count, 10)

	return err
}

// CountPostExpired ...
func (storage *partnerStorage) CountPostsExpired(partner *models.Partner) error {
	var count int64
	queryStatement := storage.Db.Model(&models.Post{}).
		Where("creator_id = ?", partner.ID).
		Count(&count)

	err := queryStatement.Error
	partner.PostExpired = strconv.FormatInt(count, 10)

	return err
}

// UpdateMentorLike ...
func (storage *partnerStorage) UpdateMentorLike(mentorID string) (result repo.RepoResponse, status repo.Status) {
	mentor, err := storage.getPartnerProfile(mentorID)
	if err != nil {
		status = repo.NotFound
		return
	}

	mentor.TotalLike++
	queryStmt := storage.Db.Model(&mentor).Updates(map[string]interface{}{"update_at": time.Now(), "total_like": mentor.TotalLike})
	err = queryStmt.Error
	if err != nil {
		status = repo.GetError
		return
	}

	// delete Password
	mentor.Password = ""

	result = repo.RepoResponse{Status: true, Data: mentor}
	status = repo.Liked
	return
}
