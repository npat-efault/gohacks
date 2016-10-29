#!/bin/sh
cat <<EOF
# gohacks [![GoDoc](https://godoc.org/github.com/npat-efault/gohacks?status.svg)](https://godoc.org/github.com/npat-efault/gohacks)

Various Go hacks.

EOF

go list -f '- **{{ .Name }}:** {{ .Doc }}

' github.com/npat-efault/gohacks/...

cat <<EOF

[Documentation at godoc.org.](https://godoc.org/github.com/npat-efault/gohacks)

EOF
