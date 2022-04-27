package books

import (
	"time"

	"github.com/owenlilly/progorm"
	"github.com/owenlilly/progorm-connection/connection"
	"gopkg.in/guregu/null.v4"
)

// Book books table model
type Book struct {
	ID          uint   `gorm:"primaryKey"`
	Author      string `gorm:"size:128"`
	Title       string `gorm:"size:256"`
	ReleaseDate *time.Time
	ISBN        null.String `gorm:"size:32"`
}

// Paged holds a page of results
type Paged struct {
	Page    uint   `json:"page"`
	PerPage uint   `json:"per_page"`
	Total   uint   `json:"total"`
	Pages   uint   `json:"pages"`
	Books   []Book `json:"books"`
}

// BookRepository repository interface for accessing books table
type BookRepository interface {
	Insert(book Book) (id uint, err error)
	FindAll(page, perPage uint) (Paged, error)
	FindByTitle(title string, page, perPage uint) (Paged, error)
}

type bookRepository struct {
	progorm.BaseRepository
}

// NewBookRepository create a new instance of BookRepository
func NewBookRepository(connMan connection.Manager) BookRepository {
	repo := bookRepository{
		BaseRepository: progorm.NewBaseRepository(connMan),
	}

	// omit this step if you're using your own migration tool
	repo.AutoMigrateOrWarn(&Book{})

	return &repo
}

// Insert insert a new book
func (r bookRepository) Insert(book Book) (id uint, err error) {
	// do some validations on book
	if err = r.InsertRecord(&book); err != nil {
		return 0, err
	}
	return book.ID, nil
}

// FindAll find books in paged results
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

// FindAll find books in paged results
func (r bookRepository) FindByTitle(title string, page, perPage uint) (Paged, error) {
	result := Paged{PerPage: perPage}

	query := r.DB().
		Model(Book{}).
		Where(Book{ISBN: null.StringFrom(title)}).
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
