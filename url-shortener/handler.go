package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func MapHandler(pathsToURLs map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToURLs[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(fallback http.Handler) (http.HandlerFunc, error) {
	file := openYAMLFile()
	defer file.Close()

	yamlLines := readFile(file)
	pathURLs := parseYAMLLines(yamlLines)

	pathsToURLs := make(map[string]string)
	for _, pathURL := range pathURLs {
		pathsToURLs[pathURL.Path] = pathURL.URL
	}

	return MapHandler(pathsToURLs, fallback), nil
}

func openYAMLFile() (file *os.File) {
	file, err := os.Open(".yaml")

	if err != nil {
		fmt.Printf("Failed to open YAML file: \n%s\n", err.Error())
		os.Exit(1)
	}

	return
}

func readFile(file *os.File) (lines []byte) {
	reader := bufio.NewReader(file)

	lines, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Failed to read file: \n%s\n", err.Error())
		os.Exit(1)
	}

	return
}

func parseYAMLLines(yamlLines []byte) []pathURL {
	var pathURLs []pathURL

	err := yaml.Unmarshal(yamlLines, &pathURLs)
	if err != nil {
		fmt.Printf("Failed to parse YAML: \n%s\n", err.Error())
		os.Exit(1)
	}

	return pathURLs
}
