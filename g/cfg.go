package g

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enabled   bool   `json:"enabled"`
	Listen    string `json:"listen"`
	WhiteList string `json:"whitelist"`
}

type SmtpConfig struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Spliter  string `json:"spliter"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type GlobalConfig struct {
	Debug         bool            `json:"debug"`
	Http          *HttpConfig     `json:"http"`
	Smtp          *SmtpConfig     `json:"smtp"`
	IgnoreMetrics map[string]bool `json:"ignore"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
