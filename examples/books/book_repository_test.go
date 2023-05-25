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

const testDbName = "test.db"

type SuiteBookRepository struct {
	suite.Suite

	connMan connection.Manager
	repo    BookRepository
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(SuiteBookRepository))
}

func (s *SuiteBookRepository) SetupSuite() {
	var err error
	// create a new SQL connection manager, there's also a postgres connection manager
	s.connMan, err = sqliteconn.NewConnectionManager(testDbName, &gorm.Config{
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

	s.repo = NewBookRepository(s.connMan)
}

func (s *SuiteBookRepository) TearDownSuite() {
	_ = os.Remove(testDbName)
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
	books := s.givenBooksExist()

	for i, book := range books {
		page, err := s.repo.FindAll(uint(i), 1)
		if s.NoError(err) {
			s.EqualValues(2, page.Pages)
			s.EqualValues(1, page.PerPage)
			s.EqualValues(i, page.Page)
			s.EqualValues(2, page.Total)
			s.EqualValues(1, len(page.Books))
			s.EqualValues(book.Title, page.Books[0].Title)
		}
	}
}

func (s *SuiteBookRepository) givenBooksExist() []Book {
	now := time.Now().UTC()
	book1 := Book{
		Author:      "George RR Martin",
		Title:       "Game Of Thrones",
		ReleaseDate: &now,
		ISBN:        null.StringFrom("978-3-18-148410-0"),
	}
	id1, err := s.repo.Insert(book1)
	s.NotEqual(0, id1)
	s.NoError(err)
	book1.ID = id1

	book2 := Book{
		Author:      "Unknown",
		Title:       "Beowulf",
		ReleaseDate: &now,
	}
	id2, err := s.repo.Insert(book2)
	s.NotEqual(0, id2)
	s.NoError(err)
	book2.ID = id2

	return []Book{book1, book2}
}
