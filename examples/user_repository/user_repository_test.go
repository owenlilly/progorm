package user_repository

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/owenlilly/progorm"
	"github.com/stretchr/testify/suite"
)

type SuiteUserRepository struct {
	suite.Suite

	connMan  progorm.ConnectionManager
	userRepo UserRepository

	email string
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteUserRepository))
}

func (s *SuiteUserRepository) SetupSuite() {
	// create a new SQL connection manager, there's also a postgres connection manager
	s.connMan = progorm.NewSQLiteConnectionManager("test.db", true)

	s.userRepo = NewUserRepository(s.connMan)

	s.email = "unit@test.com"
}

func (s *SuiteUserRepository) TearDownSuite() {
	db, _ := s.connMan.GetConnection()

	// clear all records
	db.Delete(&User{})
}

func (s SuiteUserRepository) Test0Insert() {
	user := &User{
		Email:       s.email,
		DisplayName: "Tester",
	}

	err := s.userRepo.Insert(user)

	s.NoError(err)
	s.NotEmpty(user.ID)
	s.NotEmpty(user.JoinedOn)
}

func (s SuiteUserRepository) Test1GetByEmail() {
	user, err := s.userRepo.GetByEmail(s.email)

	s.NoError(err)
	s.NotNil(user)
}
