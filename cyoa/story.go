package cyoa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Adventures []Adventure

type Adventure struct {
	Arc   string `json:"arc"`
	Story Story  `json:"story"`
}

type Story struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handler struct {
	story    Story
	template *template.Template
}

// HandlerOption Used as Functional Option to make the HtmlTemplate configurable
type HandlerOption func(h *handler)

var (
	defaultTemplate *template.Template
	htmlTemplate    = `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Choose Your Own Adventure</title>
		</head>
		<body>
			<section class="page">
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
				<p>{{.}}</p>
			{{end}}
			{{if .Options}}
				<ul>
				{{range .Options}}
				<li><a href="/{{.Arc}}">{{.Text}}</a></li>
				{{end}}
				</ul>
			{{else}}
				<h3>The End</h3>
			{{end}}
			</section>
			<style>
			body {
				font-family: helvetica, arial;
			}
			h1 {
				text-align:center;
				position:relative;
			}
			.page {
				width: 80%;
				max-width: 500px;
				margin: auto;
				margin-top: 40px;
				margin-bottom: 40px;
				padding: 80px;
				background: #FFFCF6;
				border: 1px solid #eee;
				box-shadow: 0 10px 6px -6px #777;
			}
			ul {
				border-top: 1px dotted #ccc;
				padding: 10px 0 0 0;
				-webkit-padding-start: 0;
			}
			li {
				padding-top: 10px;
			}
			a,
			a:visited {
				text-decoration: none;
				color: #6295b5;
			}
			a:active,
			a:hover {
				color: #7792a2;
			}
			p {
				text-indent: 1em;
			}
			</style>
		</body>
		</html>
	`
)

func init() {
	defaultTemplate = template.Must(template.New("").Parse(htmlTemplate))
}

// WithTemplate Creates a Closure Function to set a HTML Template to a Handler
func WithTemplate(tmpl *template.Template) HandlerOption {
	return func(h *handler) {
		h.template = tmpl
	}
}

// NewHandler Creates a http.HandlerFunc for Stories
func NewHandler(adventures Adventures, opts ...HandlerOption) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arcName := transformPathIntoArcName(r.URL.Path)

		if adventure := searchAdventureByArcName(adventures, arcName); adventure != nil {
			handler := handler{adventure.Story, defaultTemplate}
			for _, opt := range opts {
				opt(&handler)
			}

			err := handler.template.Execute(w, adventure.Story)
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			}
			return
		}

		errMsg := fmt.Sprintf("Adventure [%v] not found.", arcName)
		http.Error(w, errMsg, http.StatusNotFound)
	}
}

func transformPathIntoArcName(path string) (arcName string) {
	arcName = strings.TrimSpace(path)
	if arcName == "/" || arcName == "" {
		arcName = "/intro"
	}

	return arcName[1:]
}

func searchAdventureByArcName(adventures Adventures, arcName string) *Adventure {
	for _, a := range adventures {
		if a.Arc == arcName {
			return &a
		}
	}

	return nil
}

// JSONAdventures Decodes an JSON file into Adventures struct
func JSONAdventures(reader io.Reader) (Adventures, error) {
	decoder := json.NewDecoder(reader)

	var a Adventures
	if err := decoder.Decode(&a); err != nil {
		return nil, err
	}

	return a, nil
}
