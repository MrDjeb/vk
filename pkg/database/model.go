package database

type DBTabler interface {
	Migrate() error
	Insert(r DBTabler) error
}

type Tables struct {
	Photos DBPhotos
}

type Photo struct {
	PhotoID int
}
