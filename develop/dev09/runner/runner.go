package runner

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
)

// App represents parsed command line arguments
type App struct {
	link       *url.URL
	outputFile string
	depth      int
	recursive  bool
	resources  bool
}

func CLI(args []string) int {
	var app App
	err := app.fromArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		return 2
	}
	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

// fromArgs parses command line arguments into App struct
func (app *App) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-wget", flag.ContinueOnError)
	fl.StringVar(&app.outputFile, "O", "", "Path to output file")
	fl.IntVar(&app.depth, "l", -1, "Maximum number of links to follow when building downloading the site. By default depth is not set")
	fl.BoolVar(&app.recursive, "r", false, "Turn on recursive retriving")
	fl.BoolVar(&app.resources, "p", false, "Download all the files that are necessary to properly display a given HTML page")

	if err := fl.Parse(args); err != nil {
		return err
	}

	u, err := url.Parse(fl.Arg(0))
	if err != nil {
		return err
	}
	app.link = u
	app.depth++

	if app.outputFile == "" {
		app.outputFile = path.Base(app.link.Path)
	}

	return nil
}

func (app *App) run() error {
	if app.recursive {
		queue := []string{app.link.String()}
		if err := os.Mkdir(app.link.Host, os.ModePerm); err != nil {
			return err
		}
		sm := NewSitemap(app.link.String(), app.link.Host)
		err := sm.DownloadSite(queue, app.depth)
		if err != nil {
			return err
		}
		return nil
	}

	return downloadFile(app.link.String(), app.outputFile)
}
