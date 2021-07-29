package course

import (
	"errors"
	"strings"
	"time"

	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// courseStorage ...
type courseStorage struct {
	Db *gorm.DB
}

// NewCourseRepository ...
func NewCourseRepository(db *gorm.DB) repo.CourseRepository {
	return &courseStorage{
		Db: db,
	}
}

func (courseStorage *courseStorage) CreateCourse(course models.CourseEntity) (result models.CourseEntity, status repo.Status) {
	Find := courseStorage.Db.Create(&course)
	err := Find.Error

	if err != nil {
		status = repo.CanNotCreate
		return
	}

	status = repo.Created
	result = course

	return
}

func (courseStorage *courseStorage) UpdateCourse(course models.CourseEntity) (result models.CourseEntity, status repo.Status) {

	queryStatement := courseStorage.Db.Where("user_id = ?", course.UserID).Model(&course).Omit("id", "user_id").Updates(course)

	queryErr := queryStatement.Error

	if queryErr != nil {
		status = repo.CanNotUpdate
		return
	}

	status = repo.Updated
	result = course

	return
}

func (courseStorage *courseStorage) GetAllCourseByMentorId(mentorId string) (courses []models.CourseEntity, status repo.Status) {

	if err := courseStorage.Db.Where("user_id = ?", mentorId).Find(&courses).Error; err != nil {
		status = repo.CannotGetAll
		return
	}

	status = repo.Success
	return

}
func (courseStorage *courseStorage) DeleteCourseById(courseId string, mentorId string) (models.CourseEntity, repo.Status) {

	var courese models.CourseEntity

	if err := courseStorage.Db.Where("id = ? AND user_id = ?", courseId, mentorId).Delete(&courese).Error; err != nil {
		return courese, repo.CanNotDelete
	}

	return courese, repo.Success

}

func (courseStorage *courseStorage) GetCourseById(courseId string) (models.CourseEntity, repo.Status) {
	var courese models.CourseEntity

	if err := courseStorage.Db.Preload("Comments").Where("id = ?", courseId).First(&courese).Error; err != nil {
		return courese, repo.CannotGet
	}

	return courese, repo.Success
}

func (courseStorage *courseStorage) GetAllCourses() (courses []models.CourseEntity, status repo.Status) {

	if err := courseStorage.Db.Find(&courses).Error; err != nil {
		status = repo.CannotGetAll
		return
	}

	status = repo.Success
	return

}

func (courseStorage *courseStorage) GetCoursesFeature() (courses []models.CourseEntity, status repo.Status) {

	if err := courseStorage.Db.Limit(6).Where("feature <> ?", 0).Find(&courses).Error; err != nil {
		status = repo.CannotGetAll
		return
	}

	status = repo.Success
	return
}

func (courseStorage *courseStorage) CreateComment(userID string, course models.CourseEntity, comment models.Comment) (commentUpdated models.Comment, status repo.Status) {

	comment.UserID = userID
	comment.ContentID = course.ID

	courseQueryStmt := courseStorage.Db.Model(&course).
		Updates(map[string]interface{}{"update_at": time.Now().Unix()})

	commentQueryStmt := courseStorage.Db.Create(&comment)

	err := courseQueryStmt.Error
	if err != nil {
		status = repo.GetError
		return
	}

	err = commentQueryStmt.Error
	if err != nil {
		status = repo.GetError
		return
	}

	status = repo.Created
	commentUpdated = comment
	return
}

// GetComment ...
func (courseStorage *courseStorage) GetComment(courseID string) (comments []models.Comment, status repo.Status) {
	var course models.CourseEntity
	var err error
	course.ID, err = uuid.FromString(courseID)
	if err != nil {
		status = repo.CannotGetAll
		return
	}

	queryStmt := courseStorage.Db.
		Model(&course).
		Where("delete_at IS NULL").
		Association("Comments").
		Find(&comments)

	err = errors.New(queryStmt.Error())
	if err != nil {
		status = repo.CannotGetAll
		return
	}

	status = repo.Success
	return
}

// valid same value in an array string
func isRegistered(sample string, stringArray []string) bool {
	for _, value := range stringArray {
		if ok := strings.Compare(sample, value); ok == 0 {
			return true
		}
	}

	return false
}

// RegisterCourse ...
func (courseStorage *courseStorage) RegisterCourse(course models.CourseEntity, userID string) (status repo.Status) {
	if isRegistered(userID, course.Mentees) {
		status = repo.CanNotUpdate
		return
	}

	course.Mentees = append(course.Mentees, userID)
	course.Attended++

	queryStmt := courseStorage.Db.
		Model(&course).
		Updates(map[string]interface{}{"update_at": time.Now(), "mentees": course.Mentees, "attended": course.Attended})

	err := queryStmt.Error
	if err != nil {
		status = repo.CanNotUpdate
		return
	}

	status = repo.Updated
	return

}

// GetMentees ...
func (courseStorage *courseStorage) GetMentees(courseID string) (listMenteesID []string, status repo.Status) {
	var course models.CourseEntity

	querStmt := courseStorage.Db.
		Select("mentees").
		Where("id = ?", courseID)

	find := querStmt.First(&course)
	err := find.Error
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		status = repo.CannotGet
		return
	}

	listMenteesID = course.Mentees

	status = repo.Success
	return
}
