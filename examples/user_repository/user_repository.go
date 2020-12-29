package user_repository

import (
	"errors"
	"time"

	"github.com/owenlilly/progorm"
	"gorm.io/gorm"
)

// User users table model
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"size:128"`
	DisplayName string `gorm:"size:50"`
	JoinedOn    time.Time
}

// BeforeCreate perform some pre-insert operation
func (u *User) BeforeCreate(*gorm.DB) error {
	if u.JoinedOn.IsZero() {
		u.JoinedOn = time.Now().UTC()
	}
	// do some more validations...
	return nil
}

// UserRepository repository interface for accessing users table
type UserRepository interface {
	Insert(user *User) error
	GetByEmail(email string) (*User, error)
}

type userRepository struct {
	progorm.BaseRepository
}

// NewUserRepository create a new instance of UserRepository
func NewUserRepository(connMan progorm.ConnectionManager) UserRepository {
	r := &userRepository{BaseRepository: progorm.NewBaseRepository(connMan)}

	r.AutoMigrateOrWarn(&User{})

	return r
}

// Insert insert a new user
func (r userRepository) Insert(user *User) error {
	return r.InsertRecord(user)
}

// GetByEmail get a user by email
func (r userRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := r.DB().First(&user, User{Email: email})
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
