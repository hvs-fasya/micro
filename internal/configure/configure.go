package configure

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//Cfg package level global variable
var Cfg *Config

//Config app configuration structure
type Config struct {
	Env    string  `config:"Env"`
	Server *Server `config:"Server"`
}

//Server config http server data
type Server struct {
	Host        string      `config:"Host"`
	Port        string      `config:"Port"`
	Timeout     duration    `config:"Timeout"`
	Logger      *zap.Logger `config:"-"`
	RequestsLog string      `config:"RequestsLog"`
}

type duration struct {
	time.Duration
}

//UnmarshalText method satisfying toml unmarshal interface
func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return errors.Wrap(err, "unmarshal duration type error")
}

//LoadConfigs apply configuration in certain sources order
func LoadConfigs(dir string) error {
	Cfg = new(Config)
	Cfg.Server = &Server{}
	loader := confita.NewLoader(
		file.NewBackend(dir+"default.toml"), //load defaults
		env.NewBackend(),                    //load environments
	)
	err := loader.Load(context.Background(), Cfg)
	if err != nil {
		return errors.Wrap(err, "default configs load error") //load config dependent on env
	}
	loader = confita.NewLoader(file.NewBackend(dir + Cfg.Env + ".toml"))
	err = loader.Load(context.Background(), Cfg)
	if err != nil {
		return errors.Wrap(err, "env configs load error")
	}
	return nil
}
