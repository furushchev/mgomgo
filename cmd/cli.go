package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/furushchev/mgomgo"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"net/http"
	_ "net/http/pprof"
)

func ActionMigrate(c *cli.Context) {
	conn := int(c.Uint("concurrent"))
	timeout := time.Duration(c.Uint("timeout")) * time.Second
	if c.NArg() != 2 {
		logrus.Fatalln("arguments must be 2: [from uri] [to uri]")
	}
	if err := mgomgo.Migrate(c.Args().Get(0), c.Args().Get(1), conn, timeout); err != nil {
		logrus.Fatalln(err)
	}
}

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	// profiling
	go func() {
		logrus.Infoln(http.ListenAndServe("localhost:4000", nil))
	}()

	// logging
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logrus.SetFormatter(formatter)

	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "migrate inter mongo database"
	app.Version = "1.0.6"
	app.Author = "Yuki Furuta"
	app.Email = "furushchev@jsk.imi.i.u-tokyo.ac.jp"

	app.Action = ActionMigrate
	app.CommandNotFound = cmdNotFound
	app.ArgsUsage = "from to"

	app.Flags = []cli.Flag{
		cli.UintFlag{
			Name:  "concurrent, c",
			Value: 1,
			Usage: "Concurrent job number",
		},
		cli.UintFlag{
			Name:  "timeout, t",
			Value: 60,
			Usage: "Timeout for connection",
		},
	}

	app.Run(os.Args)
}

func cmdNotFound(c *cli.Context, command string) {
	logrus.Errorf(
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}
