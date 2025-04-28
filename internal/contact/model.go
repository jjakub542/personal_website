package contact

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type MessageForm struct {
	Id        string
	Name      string
	Email     string
	Message   string
	CreatedAt time.Time
}

func (u *MessageForm) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 64)),
		validation.Field(&u.Email, validation.Required, is.Email, validation.Length(1, 255)),
		validation.Field(&u.Message, validation.Required, validation.Length(1, 1024)),
	)
}

type Repository interface {
	CreateOne(*MessageForm) error
	GetAll() ([]MessageForm, error)
	DeleteOneById(string) error
}
