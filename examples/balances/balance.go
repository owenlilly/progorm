package balances

import (
	"database/sql"

	"github.com/owenlilly/progorm"
	"github.com/owenlilly/progorm/connection"
	"gorm.io/gorm"
)

// Balance balances table model
type Balance struct {
	ID    string `gorm:"primaryKey"`
	Total int64
}

// BalanceRepository repository interface for accessing balances table
type BalanceRepository interface {
	Add(balanceID string, amount int64) error
	GetBalance(balanceID string) (*Balance, error)
	Insert(balance *Balance) error

	Begin(opts ...*sql.TxOptions) (tx *gorm.DB)
	WithTx(tx *gorm.DB) BalanceRepository
	Commit() error
	Rollback() error
}

type balanceRepository struct {
	progorm.BaseRepository
}

// NewBalanceRepository create a new instance of BalanceRepository
func NewBalanceRepository(connMan connection.Manager) BalanceRepository {
	r := &balanceRepository{BaseRepository: progorm.NewBaseRepository(connMan)}

	r.AutoMigrateOrWarn(&Balance{})

	return r
}

func (r balanceRepository) Insert(balance *Balance) error {
	return r.InsertRecord(&balance)
}

func (r balanceRepository) Add(balanceID string, amount int64) error {
	return r.DB().
		Model(&Balance{ID: balanceID}).
		Update("total", gorm.Expr("total + ?", amount)).
		Error
}

func (r balanceRepository) GetBalance(balanceID string) (*Balance, error) {
	var bal Balance
	err := r.DB().Model(&Balance{}).First(&bal, &Balance{ID: balanceID}).Error
	if err != nil {
		return nil, err
	}
	return &bal, nil
}

// region: Transaction section

func (r *balanceRepository) Begin(opts ...*sql.TxOptions) (tx *gorm.DB) {
	return r.DB().Begin(opts...)
}

func (r *balanceRepository) WithTx(tx *gorm.DB) BalanceRepository {
	return &balanceRepository{BaseRepository: r.BaseRepository.WithTx(tx)}
}

func (r *balanceRepository) Commit() error {
	return r.BaseRepository.Commit()
}

func (r *balanceRepository) Rollback() error {
	return r.BaseRepository.Rollback()
}

// endregion: Transaction section
