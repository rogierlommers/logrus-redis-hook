package logrus_redis

import (
	"github.com/Sirupsen/logrus"
)

type MyHookConfig struct {
	Address string `json:"address"`
}

func init() {
	logrus_mate.RegisterHook("myhook", NewMyHook)
}

func NewMyHook(options logrus_mate.Options) (hook logrus.Hook, err error) {
	conf := MyHookConfig{}
	if err = options.ToObject(&conf); err != nil {
		return
	}

	// write your hook logic code here

	return
}
