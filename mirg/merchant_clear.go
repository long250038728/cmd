package mirg

import (
	"github.com/long250038728/web/tool/configurator"
	"github.com/long250038728/web/tool/persistence/es"
	"github.com/long250038728/web/tool/persistence/orm"
)

func NewDb() (db *orm.Gorm, readDb *orm.Gorm) {
	var configPath = "./config/online/db.yaml"
	var config orm.Config
	configurator.NewYaml().MustLoad(configPath, &config)
	db, err := orm.NewMySQLGorm(&config)
	if err != nil {
		panic(err)
	}

	configPath = "./config/online/db_read.yaml"
	configurator.NewYaml().MustLoad(configPath, &config)
	readDb, err = orm.NewMySQLGorm(&config)
	if err != nil {
		panic(err)
	}
	return
}

func NewEs() (db *es.ES) {
	var configPath = "./config/online/es.yaml"
	var config es.Config
	configurator.NewYaml().MustLoad(configPath, &config)
	db, err := es.NewEs(&config)
	if err != nil {
		panic(err)
	}
	return
}
