package search

import (
	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	"gorm.io/gorm"
)

// searchStorage ...
type searchStorage struct {
	Db *gorm.DB
}

// NewCourseRepository ...
func NewSearchRepository(db *gorm.DB) repo.SearchRepositoryable {
	return &searchStorage{
		Db: db,
	}
}

func (searchStorage *searchStorage) SearchUserName(userName string) (users []models.User, status repo.Status) {

	if err := searchStorage.Db.Where("user_name LIKE ?", "%"+userName+"%").Find(&users).Error; err != nil {
		status = repo.NotFound
		return
	}

	var usersAdvance []models.User

	for _, user := range users {

		if err := searchStorage.Db.Preload("Skills").Preload("Languages").Preload("OldJobs").Preload("Educations").Preload("TimeLines", "delete_at IS NULL").Where("id = ?", user.ID).First(&user).Error; err != nil {
			continue
		}

		usersAdvance = append(usersAdvance, user)
	}

	status = repo.Accepted
	users = usersAdvance

	return
}

// SearchPost ...
func (searchStorage *searchStorage) SearchPost(title string) (posts []models.Post, status repo.Status) {

	if err := searchStorage.Db.Where("title LIKE ?", "%"+title+"%").Find(&posts).Error; err != nil {
		status = repo.NotFound
		return
	}

	status = repo.Accepted

	return
}

// SearchCourse ...
func (searchStorage *searchStorage) SearchCourse(title string) (course []models.CourseEntity, status repo.Status) {

	if err := searchStorage.Db.Where("title LIKE ?", "%"+title+"%").Find(&course).Error; err != nil {
		status = repo.NotFound
		return
	}

	status = repo.Accepted

	return
}
