package validator

import (
	"github.com/rgrs-x/service/api/models"
	"gopkg.in/go-playground/validator.v9"
)

// UserValidator ...
type UserValidator interface {
	Valid(input models.User) error
	Handle(fieldErr []validator.FieldError) int
}
