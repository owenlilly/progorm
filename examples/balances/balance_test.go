package balances

import (
	"crypto/rand"
	"fmt"
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

type SuiteTransaction struct {
	suite.Suite

	connMan     connection.Manager
	balanceRepo BalanceRepository
}

func TestSuiteTransaction(t *testing.T) {
	suite.Run(t, new(SuiteTransaction))
}

func (s *SuiteTransaction) SetupSuite() {
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

	s.balanceRepo = NewBalanceRepository(s.connMan)
}

func (s *SuiteTransaction) TearDownSuite() {
	db, _ := s.connMan.GetConnection()

	// clear all records
	db.Exec("DELETE FROM balances")
}

func (s *SuiteTransaction) TestGivenBalanceExists_WhenAddAndCommit_VerifyComitted() {
	var balance = s.givenBalanceExists()

	// begin a transaction
	tx := s.balanceRepo.Begin()
	txRepo := s.balanceRepo.WithTx(tx)

	s.NoError(txRepo.Add(balance.ID, 10))

	// commit the transaction
	s.NoError(txRepo.Commit())

	foundBal, err := s.balanceRepo.GetBalance(balance.ID)

	if s.NoError(err) {
		s.NotEmpty(foundBal)
	}
}

func (s *SuiteTransaction) TestGivenBalanceExists_WhenAddAndRollback_VerifyRolledBack() {
	var balance = s.givenBalanceExists()

	// begin a transaction
	tx := s.balanceRepo.Begin()
	txRepo := s.balanceRepo.WithTx(tx)

	s.NoError(txRepo.Add(balance.ID, 10))

	// fetching updated balance outside of the transaction returns the original state
	nonTxBal, err := s.balanceRepo.GetBalance(balance.ID)
	if s.NoError(err) {
		s.NotEmpty(nonTxBal)
		s.Empty(nonTxBal.Total)
	}

	// fetching updated balance inside the transaction returns the updated state
	txBal, err := txRepo.GetBalance(balance.ID)
	if s.NoError(err) {
		s.NotEmpty(txBal)
		s.NotEmpty(txBal.Total)
	}

	// rollback the transaction
	s.NoError(txRepo.Rollback())

	// fetching the balance after the transaction rollback returns the original state
	rolledBackBal, err := s.balanceRepo.GetBalance(balance.ID)
	if s.NoError(err) {
		s.NotEmpty(rolledBackBal)
		s.Empty(rolledBackBal.Total)
	}
}

func (s *SuiteTransaction) givenBalanceExists() Balance {
	var balance = Balance{
		ID:    pseudoUUID(),
		Total: 0,
	}
	if err := s.balanceRepo.Insert(&balance); err != nil {
		panic(err)
	}
	return balance
}

func pseudoUUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
