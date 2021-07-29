package repository

import "github.com/rgrs-x/service/api/models"

// CommentRepository represent comment's repository contract
type CommentRepository interface {
	GetCommentByID(string) (models.Comment, Status)
	UpdateComment(models.Comment) (models.Comment, Status)
	DeleteComment(models.Comment) Status
}
