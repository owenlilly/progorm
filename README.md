[![go report card](https://goreportcard.com/badge/github.com/owenlilly/progorm "go report card")](https://goreportcard.com/report/github.com/owenlilly/progorm)
[![GoDoc](https://godoc.org/github.com/owenlilly/progorm?status.svg)](https://godoc.org/github.com/owenlilly/progorm)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

# Installation

Run `go get -u github.com/owenlilly/progorm`.

# Why Progorm

- Clean typed Repository Pattern (requires Go 1.18)
- Elegant transaction support
- Easy to use/extend
- Exposes all Gorm's underlying features  

# Usage

See [examples](https://github.com/owenlilly/progorm/tree/master/examples/user_repository) for usage.

```go
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
    progorm.BaseTypedRepository[User]
}

// NewUserRepository create a new instance of UserRepository
func NewUserRepository(connMan connection.Manager) UserRepository {
    r := &userRepository{BaseTypedRepository: progorm.NewBaseTypedRepository[User](connMan)}
    
    // optional
    r.AutoMigrateOrWarn(User{})
    
    return r
}

// Insert insert a new user
func (r userRepository) Insert(user *User) error {
    return r.InsertRecord(user)
}

// GetByEmail get a user by email
func (r userRepository) GetByEmail(email string) (*User, error) {
    user, err := r.First(User{Email: email})
    
    if err != nil {
        if r.IsRecordNotFoundError(err) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    return user, nil
}
```

---

MIT License
