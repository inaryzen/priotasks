package common

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	Debug             bool
	ServerPort        int
	DumpOnStartup     bool
	AutomaticDump     bool
	LoadDumpOnStartup string
}

var Conf Config

func InitConfig() {
	var debug = flag.Bool("d", false, "enable debug")
	var serverPort = flag.Int("p", 12345, "server port")
	var dumpOnStartup = flag.Bool("dump", true, "dump on startup")
	var autoDump = flag.Bool("auto-dump", true, "enable automatic dump")
	var loadDump = flag.String("load-dump", "", "load a specified dump on startup")
	flag.Parse()
	Conf = Config{
		Debug:             *debug,
		ServerPort:        *serverPort,
		DumpOnStartup:     *dumpOnStartup,
		AutomaticDump:     *autoDump,
		LoadDumpOnStartup: *loadDump,
	}
}

func IsDebug() bool {
	return Conf.Debug
}

func Debug(format string, v ...any) {
	if IsDebug() {
		log.Printf(format+"\n", v)
	}
}

func ResolveAppDir() (string, error) {
	appDir := "priotasks"

	usr, err := user.Current()
	if err != nil {
		log.Printf("%v", err)
		return "", err
	}

	appDir = filepath.Join(usr.HomeDir, appDir)
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err = os.Mkdir(appDir, 0755)
		if err != nil {
			log.Printf("%v", err)
			return "", err
		}
	}

	return appDir, nil
}
