package torrents_repo

import "time"

// Torrents
//go:generate gormgen -structs Torrents -input .
type Torrents struct {
	Id         int64     //
	InfoHash   []byte    //
	TotalSize  uint64    //
	Name       string    //
	CreateOn   time.Time `gorm:"time"` //
	ModifiedOn time.Time `gorm:"time"` //
}
