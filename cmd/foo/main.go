package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/kashav/serve_static"
	"gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s config.yaml", os.Args[0])
	}

	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	var configs []*serve_static.Config
	if err = yaml.Unmarshal(b, &configs); err != nil {
		log.Fatal(err.Error())
	}

	builders := make(map[string]*serve_static.Builder)
	for _, c := range configs {
		if err = c.Check(); err != nil {
			log.Fatal(err.Error())
		}
		builders[c.ID] = serve_static.NewBuilder(c)
		if err = builders[c.ID].Initialize(); err != nil {
			log.Fatal(err.Error())
		}
	}

	serve_static.NewRunner(builders).ListenAndServe()
}
