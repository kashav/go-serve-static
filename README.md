### foo

I wanted to see old versions of my github.io sites, but didn't want to type `git checkout` and `jekyll build` 18 times, so I made this instead.

Should work out of the box with any build command that outputs into a single directory with an `index.html` file.

```sh
$ echo <<EOF > config.yaml
- id: kshvmdn.com
  repo: git@github.com:kshvmdn/kshvmdn.github.io.git
  build: jekyll build
  serve: _site
EOF
$ go build -o foo cmd/foo/main.go
$ ./foo config.yaml
...
$ xdg-open "http://localhost:8080/?id=kshvmdn.com&rev=ed098544062a80ef2d8b03ca52f43e43194abb0a"
```
