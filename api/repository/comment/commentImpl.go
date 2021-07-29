package comment

import (
	"time"

	"github.com/rgrs-x/service/api/models"
	repo "github.com/rgrs-x/service/api/repository"
	"gorm.io/gorm"
)

// local storage
type commentStorage struct {
	DB *gorm.DB
}

// NewCommentRepository implement comment's interfaces
func NewCommentRepository(db *gorm.DB) repo.CommentRepository {
	return &commentStorage{
		DB: db,
	}
}

// GetCommentByID ...
func (storage *commentStorage) GetCommentByID(commentID string) (comment models.Comment, status repo.Status) {
	queryStmt := storage.DB.Table("comments").
		Where("id = ? AND delete_at IS NULL", commentID)

	find := queryStmt.First(&comment)
	err := find.Error
	if err != nil {
		status = repo.CannotGet
		return
	}

	status = repo.Success
	return
}

// UpdateComment ...
func (storage *commentStorage) UpdateComment(comment models.Comment) (contentUpdated models.Comment, status repo.Status) {
	queryStmt := storage.DB.Model(&comment).
		Updates(map[string]interface{}{
			"update_at": time.Now().Unix(),
			"content":   comment.Content,
		})

	err := queryStmt.Error
	if err != nil {
		status = repo.CanNotUpdate
		return
	}

	status = repo.Updated
	contentUpdated = comment
	return
}

// DeleteComment ...
func (storage *commentStorage) DeleteComment(comment models.Comment) (status repo.Status) {
	queryStmt := storage.DB.Model(&comment).
		Update("delete_at", time.Now().Unix())

	err := queryStmt.Error
	if err != nil {
		status = repo.CanNotDelete
		return
	}

	status = repo.Deleted
	return
}
