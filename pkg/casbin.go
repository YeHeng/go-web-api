package pkg

import (
	"github.com/YeHeng/gtool/platform/app"
	"github.com/YeHeng/gtool/platform/middleware"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	a, _ := gormadapter.NewAdapterByDB(middleware.Db)
	Enforcer, err := casbin.NewEnforcer("./etc/authz_model.conf", a)
	if err != nil {
		app.Logger.Errorf("init casbin err, %v", err)
	}
	err = Enforcer.LoadPolicy()
	if err != nil {
		app.Logger.Errorf("init casbin err, %v", err)
	}

	err = Enforcer.LoadModel()
	if err != nil {
		app.Logger.Errorf("init casbin err, %v", err)
	}
}
