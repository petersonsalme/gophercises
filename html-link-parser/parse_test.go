package link_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	link "github.com/petersonsalme/gophercises/html-link-parser"
)

var assertion = make(map[string][]string)

const ex1, ex2, ex3, ex4 = "html/ex1.html", "html/ex2.html", "html/ex3.html", "html/ex4.html"

func init() {
	assertion[ex1] = []string{"/other-page"}
	assertion[ex2] = []string{"https://www.twitter.com/joncalhoun", "https://github.com/gophercises"}
	assertion[ex3] = []string{"#", "/lost", "https://twitter.com/marcusolsson"}
	assertion[ex4] = []string{"/dog-cat"}
}

func TestParse(t *testing.T) {
	opts := []loadHtmlFuncOption{
		loadHtml(ex1), loadHtml(ex2), loadHtml(ex3), loadHtml(ex4),
	}

	for _, opt := range opts {
		file, err := opt(t)
		defer file.Close()
		if err != nil {
			t.Error(err.Error())
			return
		}

		r := bufio.NewReader(file)
		links, err := link.Parse(r)
		if err != nil {
			t.Error(err.Error())
			return
		}

		err = assert(file, links)
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
}

type loadHtmlFuncOption func(t *testing.T) (*os.File, error)

func loadHtml(name string) loadHtmlFuncOption {
	return func(t *testing.T) (*os.File, error) {
		file, err := os.Open(name)
		if err != nil {
			return nil, err
		}

		return file, nil
	}
}

func assert(file *os.File, links []link.Link) error {
	if expectedLinks, ok := assertion[file.Name()]; ok {

		if amountOfExpectedLinks, amountOfLinks := len(expectedLinks), len(links); amountOfExpectedLinks != amountOfLinks {
			return fmt.Errorf("[%d] links were expected, but only [%d] found.\n", amountOfExpectedLinks, amountOfLinks)
		}

		for i, expectedLink := range expectedLinks {
			if foundLink := links[i].Href; expectedLink != foundLink {
				return fmt.Errorf("Expected [%s], found [%s].\n", expectedLink, foundLink)
			}
		}

	} else {
		return fmt.Errorf("File [%s] not expected.\n", file.Name())
	}

	return nil
}
