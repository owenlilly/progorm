module github.com/owenlilly/progorm

go 1.18

require (
	github.com/owenlilly/progorm/connection v0.0.0
	gorm.io/gorm v1.23.4
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
)

replace github.com/owenlilly/progorm/connection v0.0.0 => ./connection
