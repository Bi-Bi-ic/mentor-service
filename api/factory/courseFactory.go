package factory

import (
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// CourseInfoFactory ...
type CourseInfoFactory struct{}

// Courseable ...
type Courseable struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	FromTime    int64      `json:"from_time"`
	ToTime      int64      `json:"to_time"`
	Street      string     `json:"street"`
	District    string     `json:"district"`
	Ward        string     `json:"ward"`
	City        string     `json:"city"`
	Attended    int        `json:"attended"`
	Mentor      mentorAble `json:"mentor"`
	TotalLike   int        `json:"total_like"`
}

type mentorInfoFactory struct{}

type mentorAble struct {
	ID          uuid.UUID `json:"id"`
	UserName    string    `json:"username"`
	Fullname    string    `json:"fullname"`
	MailContact string    `json:"mail_contact"`
	Avatar      string    `json:"avatar"`
	Cover       string    `json:"cover"`
}

// MenteeInfoFactory ...
type MenteeInfoFactory struct{}

// Menteeable ...
type Menteeable struct {
	ID       string `json:"id"`
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
}

// Create ...
func (factory CourseInfoFactory) Create(course models.CourseEntity, mentor Userable) Courseable {

	var mentorFactory mentorInfoFactory
	mentorAble := mentorFactory.create(mentor)

	courseable := Courseable{}

	courseable.ID = course.ID
	courseable.Title = course.Title
	courseable.Description = course.Description
	courseable.ToTime = course.ToTime
	courseable.FromTime = course.FromTime
	courseable.Street = course.Street
	courseable.District = course.District
	courseable.Ward = course.Ward
	courseable.City = course.City
	courseable.Attended = course.Attended
	courseable.Mentor = mentorAble
	courseable.TotalLike = course.TotalLike
	return courseable

}

func (factory CourseInfoFactory) CreateFromList(courseEntities []models.CourseEntity, mentor Userable) (courseables []Courseable) {

	for _, value := range courseEntities {

		courseables = append(courseables, factory.Create(value, mentor))
	}

	return

}

// Create ...
func (factory MenteeInfoFactory) Create(mentee models.User) Menteeable {

	var userFactory UserInfoFactory
	userAble := userFactory.Create(mentee)

	return Menteeable{
		ID:       userAble.ID.String(),
		FullName: userAble.Fullname,
		UserName: userAble.UserName,
		Avatar:   userAble.Avatar,
	}
}

// CreateFromList for Mentees
func (factory MenteeInfoFactory) CreateFromList(userEntities []models.User) []Menteeable {
	mentees := []Menteeable{}
	for _, value := range userEntities {
		mentees = append(mentees, factory.Create(value))
	}

	return mentees
}

func (factory mentorInfoFactory) create(user Userable) mentorAble {
	return mentorAble{
		ID:          user.ID,
		UserName:    user.UserName,
		Fullname:    user.Fullname,
		MailContact: user.MailContact,
		Avatar:      user.Avatar,
		Cover:       user.Cover,
	}
}
