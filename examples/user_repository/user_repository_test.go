package user_repository

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/owenlilly/progorm"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	s.connMan = progorm.NewSQLiteConnectionManager("test.db", &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		),
	})

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
