package viper

import (
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

type errors []error

func (e errors) Error() string {
	var sb strings.Builder
	for _, err := range e {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}
	return sb.String()
}

var vipers []*viper.Viper

func Watch(files []string) error {
	var errs errors
	for _, file := range files {
		v := viper.New()
		v.SetConfigFile(file)
		if e := v.ReadInConfig(); e != nil {
			errs = append(errs, e)
			continue
		}
		if e := viper.MergeConfigMap(v.AllSettings()); e != nil {
			errs = append(errs, e)
			continue
		}
		v.OnConfigChange(configWatcher)
		v.WatchConfig()
		vipers = append(vipers, v)
	}
	return errs
}

func configWatcher(_ fsnotify.Event) {
	for _, v := range vipers {
		_ = viper.MergeConfigMap(v.AllSettings())
	}
}
