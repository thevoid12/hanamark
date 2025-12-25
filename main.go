package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	configurables "hanamark/internal"
	logs "hanamark/logger"
	"hanamark/parser"
	"hanamark/util"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
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
		fmt.Println("No command specified.")
		fmt.Println("Run './hanamark help' for usage information")
		return
	}
	fmt.Println("hiiii homies this is hanamark")
	command := args[0]
	command = strings.ToLower(command)

	switch command {
	case "init":
		err := Init("configurables")
		if err != nil {
			log.Println("unable to init new project", err)
		}
		fmt.Println("project initialized successfully")
		fmt.Println("run ./hanamark help to learn more....")
		return

	case "help":
		printHelp()
		return
	case "build":
		l, ctx, err := setupConfig()
		if err != nil {
			return
		}
		err = util.CopyAssets(viper.GetString("filepath.mdAssetsSourcePath"), viper.GetString("filepath.mdAssetsDestPath"))
		if err != nil {
			l.Sugar().Error("copy assets files failed", err)
			return
		}
		err = parser.ParseFiles(ctx)
		if err != nil {
			l.Sugar().Error("error parsing files", err)
			return
		}

		// copy css from source static file to the static file of destination
		err = util.CopyAssets(viper.GetString("filepath.sourceStaticFiles"), filepath.Join(viper.GetString("filepath.destHtmlDir"), "static"))
		if err != nil {
			l.Sugar().Error("copy css files failed", err)
			return
		}
		l.Info("::::::::::::::::::conversion successful:::::::::::::::::::::::::::::::::::::::")
	case "serve":
		l, _, err := setupConfig()
		if err != nil {
			return
		}
		servePort := viper.GetString("servePort")
		if servePort == "" {
			servePort = "3000"
		}
		l.Info("::::::::::::::::::::starting local server at port localhost" + servePort + ":::::::::::::::::::::::::::")
		dest := viper.GetString("filepath.destHtmlDir")
		serveStaticFiles(dest, servePort)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run './hanamark help' for usage information")
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

func setupConfig() (l *zap.Logger, ctx context.Context, err error) {
	exists, err := util.DirExists("./configurables")
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	if !exists {
		err = errors.New("hanamark not properly initialized run ./hanamark init")
		log.Println(err)
		return nil, nil, err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configurables") // path to look for the config file in

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("there is a error in the path of config file", err)
			return nil, nil, err
		} else {
			// Config file was found but another error was produced
			log.Println("error laoding config file from viper", err)
			return nil, nil, err
		}
	}

	l, err = logs.InitializeLogger()
	if err != nil {
		log.Println("error initializing logger", err)
		return nil, nil, err
	}

	ctx = context.Background()
	ctx = logs.SetLoggerctx(ctx, l)
	return l, ctx, nil
}
