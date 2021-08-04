package files_repo

// Files
//go:generate gormgen -structs Files -input .
type Files struct {
	Id        int64  //
	TorrentId int64  //
	Size      int64  //
	Path      string //
	IsReadme  int32  //
	Content   string //
}
