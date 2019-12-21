package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	indexURL = "https://index.golang.org/index"
	proxyURL = "https://proxy.golang.org"
	sumURL   = "https://sum.golang.org"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "please specify a subcommand: index|proxy|sum\n")
	}

	flag.Parse()

	if flag.Arg(0) == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "index":
		err := index(flag.Args()[1:])
		if err != nil {
			log.Println(err)
		}
	case "proxy":
		err := proxy(flag.Args()[1:])
		if err != nil {
			log.Println(err)
		}
	case "sum":
		err := sum(flag.Args()[1:])
		if err != nil {
			log.Println(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func index(args []string) error {
	var since string
	var limit int
	var help bool

	fs := flag.NewFlagSet("index", flag.ExitOnError)
	fs.StringVar(&since, "since", "", "the oldest allowable timestamp in RFC3339 "+time.RFC3339)
	fs.IntVar(&limit, "limit", 0, "limit the length of list to output, 0 < x <= 2000")
	fs.BoolVar(&help, "help", false, "print help message")
	fs.BoolVar(&help, "h", false, "print help message")
	err := fs.Parse(args)
	if err != nil || help {
		fs.Usage()
		return fmt.Errorf("index: help")
	}

	vals := url.Values{}
	if since != "" {
		vals.Add("since", since)
	}
	if limit > 0 && limit <= 2000 {
		vals.Add("limit", strconv.Itoa(limit))
	}
	u, _ := url.Parse(indexURL)
	u.RawQuery = vals.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("index: get %s: %w", u.String(), err)
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("index: response %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
	return nil
}

func proxy(args []string) error {
	var save, help bool
	fs := flag.NewFlagSet("proxy", flag.ExitOnError)
	fs.BoolVar(&save, "save", false, "save zip file")
	fs.BoolVar(&help, "help", false, "print help message")
	fs.BoolVar(&help, "h", false, "print help message")
	err := fs.Parse(args)
	if err != nil || help {
		fs.Usage()
		return fmt.Errorf("proxy: help")
	}
	mod, vers, err := modvers(fs.Args())
	if err != nil {
		return fmt.Errorf("proxy: %w", err)
	}

	var out io.Writer = os.Stdout
	p := path.Join(mod, "@v")
	if vers == "" {
		p = path.Join(p, "list")
	} else if save {
		p = path.Join(p, vers) + ".zip"
		f, err := os.Create(path.Base(p))
		if err != nil {
			return fmt.Errorf("proxy: create %s: %w", path.Base(p), err)
		}
		defer f.Close()
		out = f
	} else {
		p = path.Join(p, vers) + ".info"
	}

	u, _ := url.Parse(proxyURL)
	u.Path = "/" + p

	res, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("proxy: get %s: %w", u.String(), err)
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("proxy: response %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	io.Copy(out, res.Body)

	return nil
}

func sum(args []string) error {
	var help bool
	fs := flag.NewFlagSet("proxy", flag.ExitOnError)
	fs.BoolVar(&help, "help", false, "print help message")
	fs.BoolVar(&help, "h", false, "print help message")
	err := fs.Parse(args)
	if err != nil || help {
		fs.Usage()
		return fmt.Errorf("proxy: help")
	}

	mod, vers, err := modvers(fs.Args())
	if err != nil {
		return fmt.Errorf("sum: %w", err)
	} else if vers == "" {
		return fmt.Errorf("sum: expected a version as argument 2")
	}

	u, _ := url.Parse(sumURL)
	u.Path = "/" + path.Join("lookup", mod+"@"+vers)

	res, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("sum: get %s: %w", u.String(), err)
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("sum: response %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	io.Copy(os.Stdout, res.Body)

	return nil
}

func modvers(args []string) (module, version string, err error) {
	if len(args) == 0 {
		return "", "", fmt.Errorf("modvers: no module")
	}
	args = append(args, "")
	if i := strings.LastIndex(args[0], "@"); i < 0 {
		return args[0], args[1], nil
	} else {
		return args[0][:i], args[0][i+1:], nil
	}
}
