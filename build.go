package foo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Builder struct {
	c         *Config
	root      string
	sparse    string
	checkouts map[string]string
}

func NewBuilder(c *Config) *Builder {
	var b Builder
	b.c = c
	b.checkouts = make(map[string]string)
	return &b
}

func (b Builder) Repo() string { return b.c.Repo }

func (b *Builder) InitializeFromSparse() (err error) {
	log.Printf("filling in builder from sparse checkout: %s", b.c.SparseCheckoutDir)

	b.sparse = b.c.SparseCheckoutDir

	var sparse string
	b.root, sparse = path.Split(b.sparse)

	subdirs, err := ioutil.ReadDir(b.root)
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range subdirs {
		if dir.Name() == sparse {
			continue
		}

		fullpath := path.Join(b.root, dir.Name())

		var outb bytes.Buffer
		cmd := exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = fullpath
		cmd.Stdout = &outb
		if Debug {
			cmd.Stderr = os.Stderr
		}

		if cmd.Run() == nil {
			rev := strings.TrimSuffix(outb.String(), "\n")
			log.Printf("found checkout %s at path %s", rev, fullpath)
			b.checkouts[rev] = fullpath
		}
	}

	return nil
}

func (b *Builder) Initialize() (err error) {
	if b.c.SparseCheckoutDir != "" {
		return b.InitializeFromSparse()
	}

	if b.root, err = ioutil.TempDir(os.TempDir(), b.c.ID); err != nil {
		return err
	}
	if b.sparse, err = ioutil.TempDir(b.root, ""); err != nil {
		return err
	}
	log.Printf("cloning %s", b.c.Repo)
	args := []string{"clone", "--no-checkout", "--sparse", b.c.Repo, b.sparse}
	return runcmd("git", args, ".")
}

func (b Builder) path(rev string) string {
	if p, ok := b.checkouts[rev]; ok {
		return p
	}
	return ""
}

func (b Builder) BuildPath(rev string) string {
	if p := b.path(rev); p != "" {
		return filepath.Join(p, b.c.Serve)
	}
	return ""
}

func (b *Builder) checkout(rev string) error {
	path, err := ioutil.TempDir(b.root, "")
	if err != nil {
		return err
	}

	log.Printf("checking out rev %s", rev)
	args := []string{"-C", b.sparse, "worktree", "add", "--force", path, rev}
	if err = runcmd("git", args, "."); err != nil {
		return err
	}

	b.checkouts[rev] = path
	return nil
}

func (b Builder) build(rev string) error {
	dir := b.path(rev)
	if dir == "" {
		return fmt.Errorf("can't find dir for %s", rev)
	}
	args := strings.Split(b.c.Build, " ")
	return runcmd(args[0], args[1:], dir)
}

func (b *Builder) CheckoutAndBuild(rev string) error {
	if err := b.checkout(rev); err != nil {
		return err
	}
	if err := b.build(rev); err != nil {
		return err
	}
	return nil
}

func runcmd(name string, args []string, dir string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
