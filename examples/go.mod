module github.com/owenlilly/progorm/examples

go 1.18

require (
	github.com/owenlilly/progorm v0.0.0
	github.com/owenlilly/progorm/connection v0.0.0
	github.com/owenlilly/progorm/sqlite-connection v0.0.0
	github.com/stretchr/testify v1.7.1
	gopkg.in/guregu/null.v4 v4.0.0
	gorm.io/gorm v1.23.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
	gorm.io/driver/sqlite v1.3.2 // indirect
)

replace (
	github.com/owenlilly/progorm v0.0.0 => ../
	github.com/owenlilly/progorm/connection v0.0.0 => ../connection
	github.com/owenlilly/progorm/sqlite-connection v0.0.0 => ../sqlite-connection
)
