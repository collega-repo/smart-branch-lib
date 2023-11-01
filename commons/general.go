package commons

import (
	"github.com/goccy/go-json"
	"regexp"
	"time"
)

var Configs = config{}

type config struct {
	App        app        `yaml:"app"`
	Swagger    swagger    `yaml:"swagger"`
	Log        log        `yaml:"log"`
	Datasource datasource `yaml:"datasource"`
	Core       core       `yaml:"core"`
}

type app struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Desc     string `yaml:"desc"`
	LoadRole bool   `yaml:"load-role"`
	Port     int    `yaml:"port"`
	Debug    bool   `yaml:"debug"`
	Cache    cache  `yaml:"cache"`
	Stream   stream `yaml:"stream"`
}

type stream struct {
	PubSub bool `yaml:"pubSub"`
}

type cache struct {
	Key      string `yaml:"key"`
	Response int    `yaml:"response"`
}

type swagger struct {
	Host string `yaml:"host"`
	Path struct {
		Base string `yaml:"base"`
		Docs string `yaml:"docs"`
	} `yaml:"path"`
}

type log struct {
	Path string `yaml:"path"`
}

type datasource struct {
	DB    db    `yaml:"db"`
	Redis redis `yaml:"redis"`
}

type db struct {
	Driver      string        `yaml:"driver"`
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	Database    string        `yaml:"database"`
	Schema      string        `yaml:"schema"`
	Username    string        `yaml:"username"`
	Password    string        `yaml:"password"`
	Sslmode     string        `yaml:"sslmode"`
	MaxIdle     int           `yaml:"maxIdle"`
	MaxIdleTime time.Duration `yaml:"maxIdleTime"`
	MaxOpen     int           `yaml:"maxOpen"`
	MaxLifetime time.Duration `yaml:"maxLifetime"`
	Ping        bool          `yaml:"ping"`
	Debug       bool          `yaml:"debug"`
}

type redis struct {
	Network     string        `yaml:"network"`
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	DB          int           `yaml:"db"`
	Username    string        `yaml:"username"`
	Password    string        `yaml:"password"`
	IdleMin     int           `yaml:"idleMin"`
	PoolSize    int           `yaml:"poolSize"`
	PoolTimeout time.Duration `yaml:"poolTimeout"`
	Ping        bool          `yaml:"ping"`
}

type core struct {
	Host      string `yaml:"host"`
	Username  string `yaml:"username"`
	ChannelId string `yaml:"channelId"`
	SecretKey string `yaml:"secretKey"`
	IsForward bool   `yaml:"isForward"`
	LogDebug  bool   `yaml:"logDebug"`
}

var APPL_MAP = map[string]string{
	"FE_BE":       "Aplikasi BackEnd",
	"FE_CUSTOMER": "Aplikasi Nasabah",
	"FE_TELLER":   "Aplikasi Teller",
}

func ConvertStruct(src any, dest any) error {
	jsonByte, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonByte, dest)
}

var PasswordRegex = regexp.MustCompile(`"password":\s*"([^"]+)"`)
var UsernameRegex = regexp.MustCompile(`"username"\s*:\s*"([^"]+)"`)
var UserIdRegex = regexp.MustCompile(`"userId"\s*:\s*"([^"]+)"`)
