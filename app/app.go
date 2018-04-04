package app

import "github.com/spf13/viper"

const (
	Version = ""
	Commit = ""
)

func GetVersion() string {
	return Version;
}

func GetVersionName() string {
	return viper.Get("version.name").(string);
}