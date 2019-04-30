package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
)

var pathRegex *regexp.Regexp
var rootContentPath string
var leading = []byte(`
	<html>
		<head>
			<meta charset="utf-8">
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/3.0.1/github-markdown.css" />
		</head>
		<body>
			<article class="markdown-body entry-content" style="padding: 30px;">`)
var ending = []byte(`
			</article>
		</body>
	</html>`)

func main() {
	var err error

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "1337"
	}

	fmt.Printf("Server listens on port %s\n", port)

	rootContentPath = os.Getenv("CONTENT")
	if len(rootContentPath) == 0 {
		rootContentPath = "."
	}

	rootContentPath, err = filepath.Abs(rootContentPath)
	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(rootContentPath)
	if err != nil {
		panic(err)
	} else if !stat.IsDir() {
		panic(fmt.Errorf("%s is not a directory", rootContentPath))
	}

	fmt.Printf("Content is read from %s\n", rootContentPath)

	pathRegex, err = regexp.Compile("^\\/([a-zA-Z0-9_\\-]+\\/)*([a-zA-Z0-9_\\-]+\\.md)$")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !pathRegex.MatchString(r.URL.Path) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filepaths := make([]string, 1+1+strings.Count(path, "/"))
	index := 0
	filepaths[index] = rootContentPath
	index++
	path = path[1:]
	for true {
		slashIndex := strings.Index(path, "/")
		if slashIndex == -1 {
			filepaths[index] = path
			break
		}

		filepaths[index] = path[:slashIndex+1]
		index++
		path = path[slashIndex+1:]
	}

	file := filepath.Join(filepaths...)

	stat, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprintf(os.Stderr, "%s\n", err)
	} else if stat.IsDir() {
		panic(fmt.Errorf("%s is a directory", file))
	}

	mdData, err := ioutil.ReadFile(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(leading)
	w.Write(github_flavored_markdown.Markdown(mdData))
	w.Write(ending)
}
