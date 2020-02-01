package progorm

import (
	"log"

	"github.com/jinzhu/gorm"
)

type BaseRepository struct {
	connMan ConnectionManager
	db      *gorm.DB
}

func NewBaseRepository(connMan ConnectionManager) BaseRepository {
	r := BaseRepository{
		connMan: connMan,
	}

	var err error
	r.db, err = connMan.GetConnection()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	return r
}

// Model to insert must be a pointer/reference type
func (r BaseRepository) InsertRecord(model interface{}) error {
	return r.db.Create(model).Error
}

func (r BaseRepository) AutoMigrate(models ...interface{}) error {
	return r.connMan.AutoMigrate(models...)
}

func (r BaseRepository) AutoMigrateOrWarn(models ...interface{}) {
	if err := r.connMan.AutoMigrate(models...); err != nil {
		log.Println("warning:", err.Error())
	}
}

func (r BaseRepository) ConnectionManager() ConnectionManager {
	return r.connMan
}

func (r BaseRepository) DB() *gorm.DB {
	return r.db
}
