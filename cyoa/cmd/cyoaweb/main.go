package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/petersonsalme/gophercises/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	port := flag.Int("port", 3000, "Application port")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(file)
	if err != nil {
		panic(err)
	}

	storyHandler := cyoa.NewHandler(story)

	appPort := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting application at port %s\n", appPort)
	log.Fatal(http.ListenAndServe(appPort, storyHandler))
}
