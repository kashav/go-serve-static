package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/kashav/foo"
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

	var configs []*foo.Config
	if err = yaml.Unmarshal(b, &configs); err != nil {
		log.Fatal(err.Error())
	}

	builders := make(map[string]*foo.Builder)
	for _, c := range configs {
		if err = c.Check(); err != nil {
			log.Fatal(err.Error())
		}
		builders[c.ID] = foo.NewBuilder(c)
		if err = builders[c.ID].Initialize(); err != nil {
			log.Fatal(err.Error())
		}
	}

	foo.NewRunner(builders).ListenAndServe()
}
