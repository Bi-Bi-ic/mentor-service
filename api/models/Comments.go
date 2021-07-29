package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Type of Content
const (
	PostConent    = "post"
	CourseContent = "course"
)

// Comment ...
type Comment struct {
	Base
	ContentID   uuid.UUID `json:"-"`
	ContentType string    `json:"content_type"`
	Content     string    `json:"content"`
	UserID      string    `json:"-"`
	UserName    string    `json:"-" gorm:"-"`
	Fullname    string    `json:"-" gorm:"-"`
	Avatar      string    `json:"-" gorm:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (comment *Comment) BeforeCreate(scope *gorm.DB) (err error) {
	id, err := uuid.NewV4()
	if err != nil {
		return
	}

	scope.Statement.SetColumn("ID", id)
	return
}
