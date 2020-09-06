[![go report card](https://goreportcard.com/badge/github.com/owenlilly/progorm "go report card")](https://goreportcard.com/report/github.com/owenlilly/progorm)
[![GoDoc](https://godoc.org/github.com/owenlilly/progorm?status.svg)](https://godoc.org/github.com/owenlilly/progorm)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

# Installation

Run `go get -u github.com/owenlilly/progorm`.


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

// UserRepository repository interface for accessing books table
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
    if result.RowsAffected == 0 {
        return nil, errors.New("user not found")
    }
    
    return &user, nil
}
```

---

MIT License

Copyright (c) 2020 Owen Lilly

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.