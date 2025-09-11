package cfg

import (
	"fmt"
	"github.com/spf13/viper"
)

type Cfg struct {
	PORT                   string
	DBNAME                 string
	DBUSER                 string
	DBPASS                 string
	DBHOST                 string
	DBPORT                 string
	SECRET                 string
	REFRESH_TOKEN_LIVE_DAY int64
}

func LoadConfig() Cfg {
	v := viper.New()
	v.SetEnvPrefix("SERVER")
	v.SetDefault("PORT", "3000")
	v.SetDefault("DBNAME", "postgres")
	v.SetDefault("DBUSER", "")
	v.SetDefault("DBPASS", "")
	v.SetDefault("DBHOST", "localhost")
	v.SetDefault("DBPORT", "5432")
	v.SetDefault("SECRET", "0987654321")
	v.SetDefault("REFRESH_TOKEN_LIVE_DAY", 7) // days
	v.AutomaticEnv()
	//fmt.Println("SERVER_PORT:", os.Getenv("SERVER_PORT"))
	//fmt.Println("SERVER_DBNAME:", os.Getenv("SERVER_DBNAME"))
	//fmt.Println("SERVER_DBUSER:", os.Getenv("SERVER_DBUSER"))
	//fmt.Println("SERVER_DBPASS:", os.Getenv("SERVER_DBPASS"))
	//fmt.Println("SERVER_DBHOST:", os.Getenv("SERVER_DBHOST"))
	//fmt.Println("SERVER_DBPORT:", os.Getenv("SERVER_DBPORT"))
	var cfg Cfg
	err := v.Unmarshal(&cfg)
	if err != nil {
		fmt.Println("Failed to load configuration: " + err.Error())
	}
	return cfg
}

func (c *Cfg) GetDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUSER, c.DBPASS, c.DBHOST, c.DBPORT, c.DBNAME)
}
