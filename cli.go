package archdif

import (
	"log"

	"github.com/codegangsta/cli"
)

const version = "1.0.0"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Run(args []string) {
	initApp(args)
}

func initApp(args []string) {
	var (
		format        string
		keepUntracked bool
	)

	app := cli.NewApp()
	app.Name = "archdif"
	app.Usage = "Runs git diff and archives modified files as zip or tar"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format, f",
			Value:       "zip",
			Usage:       "set archive format. Available: zip, tar",
			Destination: &format,
		},
		cli.BoolFlag{
			Name:        "untracked, u",
			Usage:       "add untrackted files to the archive",
			Destination: &keepUntracked,
		},
	}
	app.Action = func(c *cli.Context) {
		files := GitCmd(&keepUntracked)
		if len(files) == 0 {
			log.Fatal("Got no files to archive")
		}

		a := ArchiverFactory(&format)
		if a == nil {
			log.Fatalf("Unsupported archive format: (%s)", format)
		}
		a.Compress(files)
	}

	app.Run(args)
}
