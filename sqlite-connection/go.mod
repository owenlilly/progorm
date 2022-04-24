module github.com/owenlilly/progorm/sqlite-connection

go 1.18

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	gorm.io/driver/sqlite v1.3.2 // indirect
	gorm.io/gorm v1.23.4 // indirect
	github.com/owenlilly/progorm/connection v0.0.0-unpublished
)

replace github.com/owenlilly/progorm/connection v0.0.0-unpublished => ../connection

