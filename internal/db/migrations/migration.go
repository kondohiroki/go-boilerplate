package migrations

type Migration struct {
	Name string
	Up   func() error
	Down func() error
}

var Migrations []*Migration
