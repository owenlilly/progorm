package user_repository

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/owenlilly/progorm-connection/connection"
	sqliteconn "github.com/owenlilly/progorm-sqlite-connection/sqliteconnection"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SuiteUserRepository struct {
	suite.Suite

	connMan  connection.Manager
	userRepo UserRepository

	email string
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteUserRepository))
}

func (s *SuiteUserRepository) SetupSuite() {
	var err error
	// create a new SQL connection manager, there's also a postgres connection manager
	s.connMan, err = sqliteconn.NewConnectionManager("test.db", &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		),
	})

	s.NoError(err)

	s.userRepo = NewUserRepository(s.connMan)
	s.email = "unit@test.com"
}

func (s *SuiteUserRepository) TearDownTest() {
	db, _ := s.connMan.GetConnection()

	// clear all records
	db.Where(gorm.Expr("id IS NOT NULL")).Delete(&User{})
}

func (s *SuiteUserRepository) Test_Insert() {
	user, err := s.insertUser()

	s.NoError(err)
	s.NotEmpty(user.ID)
	s.False(user.JoinedOn.IsZero())
}

func (s *SuiteUserRepository) Test_GetByEmail() {

	_, err := s.insertUser()
	if !s.NoError(err) {
		return
	}

	user, err := s.userRepo.GetByEmail(s.email)

	s.NoError(err)
	s.NotNil(user)
}

func (s *SuiteUserRepository) insertUser() (*User, error) {
	user := &User{
		Email:       s.email,
		DisplayName: "Tester",
	}

	err := s.userRepo.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
