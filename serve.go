package foo

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

type Runner struct {
	builders map[string]*Builder
}

func NewRunner(bs map[string]*Builder) *Runner { return &Runner{bs} }

func (rr Runner) query(r *http.Request, param string) string {
	v := r.URL.Query().Get(param)
	if v == "" {
		// This might be a resource request, in which case the query params will be
		// part of the referer URL.
		if url, _ := url.Parse(r.Referer()); url != nil {
			v = url.Query().Get(param)
		}
	}
	return v
}

func (rr Runner) builder(id string) *Builder {
	if len(rr.builders) == 1 {
		for _, v := range rr.builders {
			return v
		}
	}
	if id == "" {
		return nil
	}
	if v, ok := rr.builders[id]; ok {
		return v
	}
	return nil
}

func (rr Runner) handler(w http.ResponseWriter, r *http.Request) {
	rev := rr.query(r, "rev")
	if rev == "" {
		// Maybe this should respond with a list of all revs instead?
		rev = "master"
	}

	b := rr.builder(rr.query(r, "id"))
	if b == nil {
		http.Error(w, "missing `id` parameter?", http.StatusBadRequest)
		return
	}

	path := b.BuildPath(rev)
	if path == "" {
		if err := b.CheckoutAndBuild(rev); err != nil {
			log.Printf(err.Error())
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		path = b.BuildPath(rev)
	}

	http.ServeFile(w, r, filepath.Join(path, r.URL.Path))
}

func (rr Runner) ListenAndServe() {
	http.HandleFunc("/", rr.handler)
	log.Printf("listening on port %d...", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}
