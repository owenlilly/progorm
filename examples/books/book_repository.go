package books

import (
	"time"

	"github.com/owenlilly/progorm"
)

// Represents `public`.`books` table in database
type Book struct {
	ID          uint   `gorm:"primary_key"`
	Author      string `gorm:"size:128"`
	Title       string `gorm:"size:256"`
	ReleaseDate *time.Time
	ISBN        string `gorm:"size:32"`
}

// FindRecords all result set
type Paged struct {
	Page    uint   `json:"page"`
	PerPage uint   `json:"per_page"`
	Total   uint   `json:"total"`
	Pages   uint   `json:"pages"`
	Books   []Book `json:"books"`
}

// Using an interface makes mocking easier
type BookRepository interface {
	Insert(book Book) (id uint, err error)
	FindAll(page, perPage uint) (Paged, error)
}

type bookRepository struct {
	progorm.BaseRepository
}

// Create a new instance of BookRepository
func NewBookRepository(connMan progorm.ConnectionManager) BookRepository {
	repo := bookRepository{
		BaseRepository: progorm.NewBaseRepository(connMan),
	}

	// omit this step if you're using your own migration tool
	repo.AutoMigrateOrWarn(&Book{})

	return &repo
}

// Insert a new book
func (r bookRepository) Insert(book Book) (id uint, err error) {
	// do some validations on book
	if err = r.InsertRecord(&book); err != nil {
		return 0, err
	}
	return book.ID, nil
}

// Find all books and page results
func (r bookRepository) FindAll(page, perPage uint) (Paged, error) {
	result := Paged{PerPage: perPage}

	query := r.DB().
		Model(Book{}).
		Order("title ASC")

	pageInfo, err := r.FindRecords(page, perPage, query, &result.Books)
	if err != nil {
		return result, err
	}

	result.Page = pageInfo.Page
	result.Total = pageInfo.Total
	result.Pages = pageInfo.Pages

	return result, nil
}
