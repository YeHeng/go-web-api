package persistence

import (
	"bytes"
	"database/sql"
	"fmt"
	repo "github.com/YeHeng/go-web-api/internal/pkg/db"
	"github.com/YeHeng/go-web-api/internal/pkg/magneticod/util"
	"github.com/YeHeng/go-web-api/pkg/errors"
	"text/template"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"
)

type mysqlDatabase struct {
}

func makeSqlite3Database() (Database, error) {
	d := new(mysqlDatabase)
	return d, nil
}

func (mysql *mysqlDatabase) DoesTorrentExist(infoHash []byte) (bool, error) {

	db := repo.GetDb()

	rows, err := db.Raw("SELECT 1 FROM torrents WHERE info_hash = ?", infoHash).Rows()
	if err != nil {
		return false, err
	}
	defer closeRows(rows)

	exists := rows.Next()
	if rows.Err() != nil {
		return false, err
	}

	return exists, nil
}

func (mysql *mysqlDatabase) AddNewTorrent(infoHash []byte, name string, files []File) error {

	db := repo.GetDb()
	if !utf8.ValidString(name) {
		zap.L().Warn(
			"Ignoring a torrent whose name is not UTF-8 compliant.",
			zap.ByteString("infoHash", infoHash),
			zap.Binary("name", []byte(name)),
		)

		return nil
	}

	db.Begin()
	err := db.Begin().Error
	if err != nil {
		return errors.Wrap(err, "conn.Begin")
	}
	// If everything goes as planned and no error occurs, we will commit the transaction before
	// returning from the function so the tx.Rollback() call will fail, trying to rollback a
	// committed transaction. BUT, if an error occurs, we'll get our transaction rollback'ed, which
	// is nice.
	defer db.Rollback()

	var totalSize uint64 = 0
	for _, file := range files {
		totalSize += uint64(file.Size)
	}

	// This is a workaround for a bug: the database will not accept total_size to be zero.
	if totalSize == 0 {
		zap.L().Debug("Ignoring a torrent whose total size is zero.")
		return nil
	}

	if exist, err := mysql.DoesTorrentExist(infoHash); exist || err != nil {
		return err
	}

	var lastInsertId int64

	err = db.Exec(`
		INSERT INTO torrents (
			info_hash,
			name,
			total_size,
			discovered_on
		) VALUES (?, ?, ?, ?)
	`, infoHash, name, totalSize, time.Now().Unix()).Error
	if err != nil {
		return errors.Wrap(err, "tx.QueryRow (INSERT INTO torrents)")
	}

	for _, file := range files {
		if !utf8.ValidString(file.Path) {
			zap.L().Warn(
				"Ignoring a file whose path is not UTF-8 compliant.",
				zap.Binary("path", []byte(file.Path)),
			)

			// Returning nil so deferred tx.Rollback() will be called and transaction will be canceled.
			return nil
		}

		err = db.Exec("INSERT INTO files (torrent_id, size, path) VALUES (?, ?, ?)",
			lastInsertId, file.Size, file.Path,
		).Error
		if err != nil {
			return errors.Wrap(err, "tx.Exec (INSERT INTO files)")
		}
	}

	err = db.Commit().Error
	if err != nil {
		return errors.Wrap(err, "tx.Commit")
	}

	return nil
}

func (mysql *mysqlDatabase) Close() error {
	return nil
}

func (mysql *mysqlDatabase) GetNumberOfTorrents() (int64, error) {

	db := repo.GetDb()

	rows, err := db.Raw("SELECT MAX(ROWID) FROM torrents").Rows()
	if err != nil {
		return 0, err
	}
	defer closeRows(rows)

	if !rows.Next() {
		return 0, fmt.Errorf("no rows returned from `SELECT MAX(ROWID)`")
	}

	var n *uint
	if err = rows.Scan(&n); err != nil {
		return 0, err
	}

	// If the database is empty (i.e. 0 entries in 'torrents') then the query will return nil.
	if n == nil {
		return 0, nil
	} else {
		return int64(*n), nil
	}
}

func (mysql *mysqlDatabase) QueryTorrents(
	query string,
	ascending bool,
	orderBy OrderingCriteria,
	limit uint,
	lastID *uint64,
	lastOrderedValue *float64,
	epoch int64,
) ([]TorrentMetadata, error) {

	db := repo.GetDb()

	if query == "" && orderBy == ByRelevance {
		return nil, fmt.Errorf("torrents cannot be ordered by relevance when the query is empty")
	}
	if (lastOrderedValue == nil) != (lastID == nil) {
		return nil, fmt.Errorf("lastOrderedValue and lastID should be supplied together, if supplied")
	}

	doJoin := query != ""
	firstPage := lastID == nil

	// executeTemplate is used to prepare the SQL query, WITH PLACEHOLDERS FOR USER INPUT.
	sqlQuery := executeTemplate(`
		SELECT id 
             , info_hash
			 , name
			 , total_size
			 , discovered_on
			 , (SELECT COUNT(*) FROM files WHERE torrents.id = files.torrent_id) AS n_files
	{{ if .DoJoin }}
			 , idx.rank
	{{ else }}
			 , 0
	{{ end }}
		FROM torrents
	{{ if .DoJoin }}
		INNER JOIN (
			SELECT rowid AS id
				 , bm25(torrents_idx) AS rank
			FROM torrents_idx
			WHERE torrents_idx MATCH ?
		) AS idx USING(id)
	{{ end }}
		WHERE     modified_on <= ?
	{{ if not .FirstPage }}
			  AND ( {{.OrderOn}}, id ) {{GTEorLTE .Ascending}} (?, ?) -- https://www.sqlite.org/rowvalue.html#row_value_comparisons
	{{ end }}
		ORDER BY {{.OrderOn}} {{AscOrDesc .Ascending}}, id {{AscOrDesc .Ascending}}
		LIMIT ?;	
	`, struct {
		DoJoin    bool
		FirstPage bool
		OrderOn   string
		Ascending bool
	}{
		DoJoin:    doJoin,
		FirstPage: firstPage,
		OrderOn:   orderOn(orderBy),
		Ascending: ascending,
	}, template.FuncMap{
		"GTEorLTE": func(ascending bool) string {
			if ascending {
				return ">"
			} else {
				return "<"
			}
		},
		"AscOrDesc": func(ascending bool) string {
			if ascending {
				return "ASC"
			} else {
				return "DESC"
			}
		},
	})

	// Prepare query
	queryArgs := make([]interface{}, 0)
	if doJoin {
		queryArgs = append(queryArgs, query)
	}
	queryArgs = append(queryArgs, epoch)
	if !firstPage {
		queryArgs = append(queryArgs, lastOrderedValue)
		queryArgs = append(queryArgs, lastID)
	}
	queryArgs = append(queryArgs, limit)

	rows, err := db.Raw(sqlQuery, queryArgs...).Rows()
	defer closeRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	torrents := make([]TorrentMetadata, 0)
	for rows.Next() {
		var torrent TorrentMetadata
		err = rows.Scan(
			&torrent.ID,
			&torrent.InfoHash,
			&torrent.Name,
			&torrent.Size,
			&torrent.DiscoveredOn,
			&torrent.NFiles,
			&torrent.Relevance,
		)
		if err != nil {
			return nil, err
		}
		torrents = append(torrents, torrent)
	}

	return torrents, nil
}

func orderOn(orderBy OrderingCriteria) string {
	switch orderBy {
	case ByRelevance:
		return "idx.rank"

	case ByTotalSize:
		return "total_size"

	case ByDiscoveredOn:
		return "discovered_on"

	case ByNFiles:
		return "n_files"

	default:
		panic(fmt.Sprintf("unknown orderBy: %v", orderBy))
	}
}

func (mysql *mysqlDatabase) GetTorrent(infoHash []byte) (*TorrentMetadata, error) {
	d := repo.GetDb()
	rows, err := d.Raw("SELECT info_hash,name,total_size,discovered_on,"+
		"(SELECT COUNT(*) FROM files WHERE torrent_id = torrents.id) AS n_files"+
		"FROM torrents WHERE info_hash = ?", infoHash).Rows()

	defer closeRows(rows)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	var tm TorrentMetadata
	if err = rows.Scan(&tm.InfoHash, &tm.Name, &tm.Size, &tm.DiscoveredOn, &tm.NFiles); err != nil {
		return nil, err
	}

	return &tm, nil
}

func (mysql *mysqlDatabase) GetFiles(infoHash []byte) ([]File, error) {
	d := repo.GetDb()
	rows, err := d.Raw(
		"SELECT size, path FROM files, torrents WHERE files.torrent_id = torrents.id AND torrents.info_hash = ?",
		infoHash).Rows()
	defer closeRows(rows)
	if err != nil {
		return nil, err
	}

	var files []File
	for rows.Next() {
		var file File
		if err = rows.Scan(&file.Size, &file.Path); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (mysql *mysqlDatabase) GetStatistics(from string, n uint) (*Statistics, error) {

	db := repo.GetDb()
	fromTime, gran, err := util.ParseISO8601(from)
	if err != nil {
		return nil, errors.Wrap(err, "parsing ISO8601 error")
	}

	var toTime time.Time
	var timef string // time format: https://www.sqlite.org/lang_datefunc.html

	switch gran {
	case util.Year:
		toTime = fromTime.AddDate(int(n), 0, 0)
		timef = "%Y"
	case util.Month:
		toTime = fromTime.AddDate(0, int(n), 0)
		timef = "%Y-%m"
	case util.Week:
		toTime = fromTime.AddDate(0, 0, int(n)*7)
		timef = "%Y-%W"
	case util.Day:
		toTime = fromTime.AddDate(0, 0, int(n))
		timef = "%Y-%m-%d"
	case util.Hour:
		toTime = fromTime.Add(time.Duration(n) * time.Hour)
		timef = "%Y-%m-%dT%H"
	}

	// TODO: make it faster!
	rows, err := db.Raw(fmt.Sprintf(`
			SELECT strftime('%s', discovered_on, 'unixepoch') AS dT
                 , sum(files.size) AS tS
                 , count(DISTINCT torrents.id) AS nD              
                 , count(DISTINCT files.id) AS nF
			FROM torrents, files
 			WHERE     torrents.id = files.torrent_id
                  AND discovered_on >= ?
                  AND discovered_on <= ?
			GROUP BY dt`,
		timef),
		fromTime.Unix(), toTime.Unix()).Rows()
	defer closeRows(rows)
	if err != nil {
		return nil, err
	}

	stats := NewStatistics()

	for rows.Next() {
		var dT string
		var tS, nD, nF uint64
		if err := rows.Scan(&dT, &tS, &nD, &nF); err != nil {
			if err := rows.Close(); err != nil {
				panic(err.Error())
			}
			return nil, err
		}
		stats.NDiscovered[dT] = nD
		stats.TotalSize[dT] = tS
		stats.NFiles[dT] = nF
	}

	return stats, nil
}

func (mysql *mysqlDatabase) setupDatabase() error {
	return nil
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		zap.L().Error("could not close row", zap.Error(err))
	}
}

func executeTemplate(text string, data interface{}, funcs template.FuncMap) string {
	t := template.Must(template.New("anon").Funcs(funcs).Parse(text))

	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		panic(err.Error())
	}
	return buf.String()
}
