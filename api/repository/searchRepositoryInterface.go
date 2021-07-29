package repository

import (
	"github.com/rgrs-x/service/api/models"
)

//UserRepository is an interface can be implemented
type SearchRepositoryable interface {
	SearchUserName(userName string) ([]models.User, Status)
	SearchPost(title string) ([]models.Post, Status)
	SearchCourse(title string) ([]models.CourseEntity, Status)
}