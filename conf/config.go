package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

var Conf Config

type Config struct {
	IsDebug   bool
	Mysql     Mysql
	OpenCheck bool
}

type Mysql struct {
	DataSource      string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

var CfgFile string

func InitConfig() {
	if _, err := toml.DecodeFile(CfgFile, &Conf); err != nil {
		panic(err)
	}
	if Conf.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
