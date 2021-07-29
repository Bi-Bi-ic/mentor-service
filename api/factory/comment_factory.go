package factory

import (
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// CommentAble ...
type CommentAble struct {
	ID          uuid.UUID `json:"id"`
	ContentType string    `json:"content_type"`
	CreateAt    int64     `json:"create_at"`
	UpdateAt    int64     `json:"update_at"`
	DeleteAt    int64     `json:"delete_at,omitempty"`
	Content     string    `json:"content"`
	User        commentor `json:"user"`
}

type commentor struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Fullname string `json:"fullname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// CommentFactory ...
type CommentFactory struct{}

// Create ...
func (factory CommentFactory) Create(comment models.Comment, creator Userable) CommentAble {
	return CommentAble{
		ID:          comment.ID,
		ContentType: comment.ContentType,
		CreateAt:    comment.CreateAt,
		UpdateAt:    comment.UpdateAt,
		DeleteAt:    comment.DeleteAt.Time.Unix(),
		Content:     comment.Content,
		User: commentor{
			ID:       comment.UserID,
			UserName: creator.UserName,
			Fullname: creator.Fullname,
			Avatar:   creator.Avatar,
		},
	}
}
