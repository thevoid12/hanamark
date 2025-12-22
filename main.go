package main

import (
	"context"
	"flag"
	"fmt"
	configurables "hanamark/internal"
	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Usage = func() {
		fmt.Println("hello world")
	}
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		return
	}
	fmt.Println("hiiii homies this is hanamark")
	command := args[0]
	command = strings.ToLower(command)

	if command == "init" {
		err := Init("configurables")
		if err != nil {
			log.Println("unable to init new project", err)
		}
		return
	}

	if command == "help" {
		printHelp()
		return
	}

	exists, err := util.DirExists("./configurables")
	if err != nil {
		log.Println(err)
		return
	}

	if !exists {
		log.Println("hanamark not properly initialized run ./hanamark init")
		return
	}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configurables") // path to look for the config file in

	err = viper.ReadInConfig()
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
	if command == "build" {

		// _, err = template.ParseGlob(filepath.Join(viper.GetString("filepath.templatePath"), "*.html"))
		// if err != nil {
		// 	l.Sugar().Error("parse glob added failed", err)
		// 	return
		// }

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
		l.Info("::::::::::::::::::conversion successful:::::::::::::::::::::::::::::::::::::::")
	} else if command == "serve" {
		servePort := viper.GetString("servePort")
		if servePort == "" {
			servePort = "3000"
		}
		l.Info("::::::::::::::::::::starting local server at port localhost" + servePort + ":::::::::::::::::::::::::::")
		dest := viper.GetString("filepath.destHtmlRoot")
		serveStaticFiles(dest, servePort)
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'hanamark help' for usage information")
	}
}

func Init(dst string) error {
	if err := util.EnsureEmptyDir(dst); err != nil {
		return err
	}

	return util.CopyEmbedDir(configurables.FS, "configurables", dst)
}

func printHelp() {
	fmt.Println("Hanamark - A static site generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  hanamark <command>")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  init     Initialize a new hanamark project in the current directory")
	fmt.Println("           Creates a 'configurables' directory with default templates and config")
	fmt.Println()
	fmt.Println("  build    Build the static site from markdown files")
	fmt.Println("           Converts markdown to HTML using templates and copies assets")
	fmt.Println("           Requires: ./configurables directory with config.json")
	fmt.Println()
	fmt.Println("  serve    Start a local development server")
	fmt.Println("           Serves the generated site from the output directory")
	fmt.Println("           Default port: 3000 (configurable via 'servePort' in config.json)")
	fmt.Println()
	fmt.Println("  help     Display this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./hanamark init              # Initialize new project")
	fmt.Println("  ./hanamark build             # Build the site")
	fmt.Println("  ./hanamark serve             # Start local server")
	fmt.Println()
}

func serveStaticFiles(dir string, port string) {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
