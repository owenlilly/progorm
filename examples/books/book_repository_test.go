package books

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/owenlilly/progorm-connection/connection"
	sqliteconn "github.com/owenlilly/progorm-sqlite-connection/sqliteconnection"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SuiteBookRepository struct {
	suite.Suite

	connMan connection.Manager
	repo    BookRepository
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteBookRepository))
}

func (s *SuiteBookRepository) SetupSuite() {
	// create a new SQL connection manager, there's also a postgres connection manager
	s.connMan = sqliteconn.NewConnectionManager("test.db", &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		),
	})

	s.repo = NewBookRepository(s.connMan)
}

func (s *SuiteBookRepository) TearDownSuite() {
	db, _ := s.connMan.GetConnection()

	// clear all records
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Book{})
}

func (s *SuiteBookRepository) TestInsert() {
	now := time.Now().UTC()
	id1, err := s.repo.Insert(Book{
		Title:       "Game Of Thrones",
		ReleaseDate: &now,
		ISBN:        null.StringFrom("978-3-18-148410-0"),
	})
	s.NotEqual(0, id1)
	s.NoError(err)

	id2, err := s.repo.Insert(Book{
		Title:       "Beowulf",
		ReleaseDate: &now,
		ISBN:        null.StringFrom("978-3-16-148410-0"),
	})
	s.NotEqual(0, id2)
	s.NoError(err)
}

func (s *SuiteBookRepository) TestFindAll() {
	s.givenBooksExist()

	page, err := s.repo.FindAll(1, 10)
	s.NoError(err)
	s.EqualValues(1, page.Pages)
	s.EqualValues(10, page.PerPage)
	s.EqualValues(1, page.Page)
	s.EqualValues(2, page.Total)
	s.EqualValues(2, len(page.Books))
	s.EqualValues("Beowulf", page.Books[0].Title)
	s.EqualValues("Game Of Thrones", page.Books[1].Title)
}

func (s *SuiteBookRepository) givenBooksExist() {
	now := time.Now().UTC()
	id1, err := s.repo.Insert(Book{
		Author:      "George RR Martin",
		Title:       "Game Of Thrones",
		ReleaseDate: &now,
		ISBN:        null.StringFrom("978-3-18-148410-0"),
	})
	s.NotEqual(0, id1)
	s.NoError(err)

	id2, err := s.repo.Insert(Book{
		Author:      "Unknown",
		Title:       "Beowulf",
		ReleaseDate: &now,
	})
	s.NotEqual(0, id2)
	s.NoError(err)
}
