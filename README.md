```sh
$ cat <<EOF > config.yaml
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
