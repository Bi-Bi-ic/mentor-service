package models

import "github.com/lib/pq"

// CourseEntity ..
type CourseEntity struct {
	Base
	Title       string         `json:"title"`
	Description string         `json:"description"`
	FromTime    int64          `json:"from_time"`
	ToTime      int64          `json:"to_time"`
	Street      string         `json:"street"`
	District    string         `json:"district"`
	Ward        string         `json:"ward"`
	City        string         `json:"city"`
	Attended    int            `json:"attended"`
	UserID      string         `json:"mentor_id"`
	Comments    []Comment      `json:"comments" gorm:"foreignkey:ContentID"`
	Mentees     pq.StringArray `json:"-" gorm:"type:text[]"`
	Feature     int            `json:"feature"`
	TotalLike   int            `json:"total_like"`
}
