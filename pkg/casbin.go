package pkg

import (
	db2 "github.com/YeHeng/go-web-api/pkg/db"
	"github.com/YeHeng/go-web-api/pkg/logger"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	a, _ := gormadapter.NewAdapterByDB(db2.Db)
	Enforcer, err := casbin.NewEnforcer("./etc/authz_model.conf", a)
	if err != nil {
		logger.Logger.Errorf("init casbin err, %v", err)
	}
	err = Enforcer.LoadPolicy()
	if err != nil {
		logger.Logger.Errorf("init casbin err, %v", err)
	}

	err = Enforcer.LoadModel()
	if err != nil {
		logger.Logger.Errorf("init casbin err, %v", err)
	}
}
