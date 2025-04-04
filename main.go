package main

import (
	"context"
	"fmt"
	"path/filepath"
	"text/template"

	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"log"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("hiiii bitch this is hanamark")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./") // path to look for the config file in

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("there is a error in the path of config file", err)
		} else {
			// Config file was found but another error was produced
			log.Println("error laoding config file from viper", err)
		}
	}

	l, err := logs.InitializeLogger()
	if err != nil {
		log.Println("error initializing logger", err)
	}

	ctx := context.Background()
	ctx = logs.SetLoggerctx(ctx, l)

	_, err = template.ParseGlob(filepath.Join(viper.GetString("filepath.templatePath"), "*.html"))
	if err != nil {
		l.Sugar().Error("parse glob added failed", err)
		return
	}

	err = parser.ParseFiles(ctx)
	if err != nil {
		l.Sugar().Error("error parsing files", err)
		return
	}
	err = util.CopyAssets(viper.GetString("filepath.sourceAssetsPath"), viper.GetString("filepath.destAssetsPath"))
	if err != nil {
		l.Sugar().Error("copy assets files failed", err)
		return
	}

	// copy css from hanamark template to dest output css template
	err = util.CopyAssets(viper.GetString("filepath.hanamarkCssPath"), viper.GetString("filepath.destCssPath"))
	if err != nil {
		l.Sugar().Error("copy css files failed", err)
		return
	}

}
