package repository

import (
	"github.com/rgrs-x/service/api/models"
)

// CourseRepository ...
type CourseRepository interface {
	CreateCourse(course models.CourseEntity) (models.CourseEntity, Status)
	UpdateCourse(course models.CourseEntity) (models.CourseEntity, Status)
	GetCourseById(courseId string) (models.CourseEntity, Status)
	GetAllCourseByMentorId(mentorId string) ([]models.CourseEntity, Status)
	DeleteCourseById(courseId string, mentorId string) (models.CourseEntity, Status)

	GetAllCourses() (courses []models.CourseEntity, status Status)

	GetCoursesFeature() (courses []models.CourseEntity, status Status)

	CreateComment(string, models.CourseEntity, models.Comment) (models.Comment, Status)
	GetComment(string) ([]models.Comment, Status)

	RegisterCourse(models.CourseEntity, string) Status
	GetMentees(string) ([]string, Status)
}
