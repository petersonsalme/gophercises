package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Story map[string]chapter

type chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []option `json:"options"`
}

type option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type handler struct {
	s Story
	t *template.Template
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
				<li><a href="/{{.Chapter}}">{{.Text}}</a></li>
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
		h.t = tmpl
	}
}

// NewHandler Creates a http.HandlerFunc for Stories
func NewHandler(s Story, opts ...HandlerOption) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)
		if path == "/" || path == "" {
			path = "/intro"
		}

		path = path[1:]

		if chapter, ok := s[path]; ok {
			handler := handler{s, defaultTemplate}
			for _, opt := range opts {
				opt(&handler)
			}

			err := handler.t.Execute(w, chapter)
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "Chapter not found...", http.StatusNotFound)
	}
}

// JSONStory Decodes an JSON file into Story struct
func JSONStory(reader io.Reader) (Story, error) {
	decoder := json.NewDecoder(reader)

	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}
