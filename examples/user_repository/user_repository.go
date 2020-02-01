package user_repository

import (
	"errors"
	"time"

	"github.com/owenlilly/progorm"
)

type User struct {
	ID          uint   `gorm:"primary_key"`
	Email       string `gorm:"size:128"`
	DisplayName string `gorm:"size:50"`
	JoinedOn    time.Time
}

// perform some pre-insert operation
func (u *User) BeforeCreate() error {
	if u.JoinedOn.IsZero() {
		u.JoinedOn = time.Now().UTC()
	}
	// do some more validations...
	return nil
}

// Always prefer interfaces when possible
type UserRepository interface {
	Insert(user *User) error
	GetByEmail(email string) (*User, error)
}

type userRepository struct {
	progorm.BaseRepository
}

func NewUserRepository(connMan progorm.ConnectionManager) UserRepository {
	r := &userRepository{BaseRepository: progorm.NewBaseRepository(connMan)}

	r.AutoMigrateOrWarn(&User{})

	return r
}

func (r userRepository) Insert(user *User) error {
	return r.InsertRecord(user)
}

func (r userRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := r.DB().First(&user, User{Email: email})
	if result.RecordNotFound() {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
