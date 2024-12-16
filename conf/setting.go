// setting - 2024/12/16
// Author: wangzx
// Description:

package conf

type Setting struct {
	App App `mapstructure:"app"          validate:"required"`
}

type App struct {
	Addr    string `yaml:"addr" validate:"required"`
	AppName string `yaml:"appname" validate:"required"`
	Mode    string `yaml:"mode" validate:"required"`
	Env     string `yaml:"env" validate:"required"`
	Token   string `yaml:"token" validate:"required"`
}
