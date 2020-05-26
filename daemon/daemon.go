package daemon

import (
	"fmt"
	"os"
	"time"

	"github.com/lj-team/go-generic/text/kv"
	"github.com/sevlyar/go-daemon"
)

func Run(conf *Config) {

	var cntxt *daemon.Context

	cntxt = &daemon.Context{
		PidFileName: conf.PidFile,
		PidFilePerm: 0644,
		LogFileName: conf.LogFile,
		LogFilePerm: 0640,
		WorkDir:     conf.WorkDir,
		Umask:       027,
		Args:        os.Args,
	}

	child, err := cntxt.Reborn()

	if err != nil {
		fmt.Println(err)
	}

	if child != nil {
		time.Sleep(time.Second)
		os.Exit(0)
	}
}

func Fork(dsn string) {
	params, _ := kv.New(dsn)

	cfg := &Config{
		PidFile: params.GetString("pid", ""),
		WorkDir: params.GetString("dir", "."),
		LogFile: params.GetString("log", ""),
	}

	Run(cfg)
}
