package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// Config struct
type Config struct {
	AddrDev      string
	AddrUser     string
	AddrWeb      string
	WebRedirURL  string
	WebPort      int
	SslCert      string
	SslKey       string
	SslCacert    string // mTLS for device
	HTTPUsername string
	HTTPPassword string
	Token        string
	FontSize     int
	WhiteList    map[string]bool
	PushToken    string
	PushTopic    string
	
}

func getConfigOpt(yamlCfg *yaml.File, name string, opt interface{}) {
	val, err := yamlCfg.Get(name)
	if err != nil {
		return
	}

	switch opt := opt.(type) {
	case *string:
		*opt = val
	case *int:
		*opt, _ = strconv.Atoi(val)
	}
}

// Parse config
func Parse(c *cli.Context) *Config {
	cfg := &Config{
		AddrDev:      c.String("addr-dev"),
		AddrUser:     c.String("addr-user"),
		AddrWeb:      c.String("addr-web"),
		WebRedirURL:  c.String("web-redir-url"),
		SslCert:      c.String("ssl-cert"),
		SslKey:       c.String("ssl-key"),
		SslCacert:    c.String("ssl-cacert"),
		HTTPUsername: c.String("http-username"),
		HTTPPassword: c.String("http-password"),
		Token:        c.String("token"),
		PushToken:    c.String("push-token"),
		PushTopic:    c.String("push-topic"),
	}

	cfg.WhiteList = make(map[string]bool)

	whiteList := c.String("white-list")

	if whiteList == "*" {
		cfg.WhiteList = nil
	} else {
		for _, id := range strings.Fields(whiteList) {
			cfg.WhiteList[id] = true
		}
	}

	yamlCfg, err := yaml.ReadFile(c.String("conf"))
	if err == nil {
		getConfigOpt(yamlCfg, "addr-dev", &cfg.AddrDev)
		getConfigOpt(yamlCfg, "addr-user", &cfg.AddrUser)
		getConfigOpt(yamlCfg, "addr-web", &cfg.AddrWeb)
		getConfigOpt(yamlCfg, "web-redir-url", &cfg.WebRedirURL)
		getConfigOpt(yamlCfg, "ssl-cert", &cfg.SslCert)
		getConfigOpt(yamlCfg, "ssl-key", &cfg.SslKey)
		getConfigOpt(yamlCfg, "ssl-cacert", &cfg.SslCacert)
		getConfigOpt(yamlCfg, "http-username", &cfg.HTTPUsername)
		getConfigOpt(yamlCfg, "http-password", &cfg.HTTPPassword)
		getConfigOpt(yamlCfg, "token", &cfg.Token)
		getConfigOpt(yamlCfg, "font-size", &cfg.FontSize)
		getConfigOpt(yamlCfg, "push-token", &cfg.PushToken)
		getConfigOpt(yamlCfg, "push-topic", &cfg.PushTopic)

		val, err := yamlCfg.Get("white-list")
		if err == nil {
			if val == "*" || val == "\"*\"" {
				cfg.WhiteList = nil
			} else {
				for _, id := range strings.Fields(val) {
					cfg.WhiteList[id] = true
				}
			}
		}
	}

	if cfg.FontSize == 0 {
		cfg.FontSize = 16
	}

	if cfg.FontSize < 12 {
		cfg.FontSize = 12
	}

	if cfg.SslCert != "" && cfg.SslKey != "" {
		_, err := os.Lstat(cfg.SslCert)
		if err != nil {
			log.Error().Msg(err.Error())
			cfg.SslCert = ""
		}

		_, err = os.Lstat(cfg.SslKey)
		if err != nil {
			log.Error().Msg(err.Error())
			cfg.SslKey = ""
		}
	}

	return cfg
}
