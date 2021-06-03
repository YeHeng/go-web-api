package middleware

import (
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/config"
	logger2 "github.com/YeHeng/go-web-api/pkg/logger"
	"os"
	"strings"
	"time"

	util2 "github.com/YeHeng/go-web-api/pkg/util"
	gorm2 "github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitDb() {

	dbConfig := config.Config.DbConfig
	userHome, _ := util2.Home()
	var err error

	err = os.MkdirAll(userHome+"/."+config.Config.AppName, 0777)

	if err != nil {
		logger2.Logger.Errorw(err.Error())
	}

	if strings.ToUpper(dbConfig.DbType) == "SQLITE" {
		Db, err = gorm.Open(sqlite.Open(userHome+"/."+config.Config.AppName+"/"+dbConfig.Dsn), &gorm.Config{
			SkipDefaultTransaction: dbConfig.SkipTransaction,
		})
	} else if strings.ToUpper(dbConfig.DbType) == "MYSQL" {
		Db, err = gorm.Open(mysql.Open(dbConfig.Dsn), &gorm.Config{})
	}
	if err != nil {
		logger2.Logger.Errorf("%v", err)
		panic(err)
	}

	Db.Config.Logger = logger.New(&optionalLogger{}, logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
		Colorful:      false,
	})
}

type optionalLogger struct {
}

func (z *optionalLogger) Printf(s string, params ...interface{}) {
	l := createLog(params)
	logger2.Logger.With(
		zap.Time("occurredAt", l.occurredAt),
		zap.String("source", l.source),
		zap.String("duration", l.duration+"ms"),
		zap.String("affectedRow", l.affectedRow),
		zap.String("message", l.message),
		zap.String("sql", l.sql),
	).Infow("")
}

type tracelog struct {
	occurredAt  time.Time
	source      string
	duration    string
	sql         string
	affectedRow string
	message     string
}

func createLog(values []interface{}) *tracelog {
	ret := &tracelog{}
	ret.occurredAt = gorm2.NowFunc()

	lens := len(values)
	if lens > 1 {
		ret.source = fmt.Sprint(values[0])

		if lens > 4 {
			ret.message = fmt.Sprint(values[1])
			ret.duration = fmt.Sprint(values[2])
			ret.affectedRow = fmt.Sprint(values[3])
			ret.sql = fmt.Sprint(values[4])
		} else if lens > 2 {
			ret.duration = fmt.Sprint(values[1])
			ret.affectedRow = fmt.Sprint(values[2])
			ret.sql = fmt.Sprint(values[3])
		}
	}
	return ret
}
