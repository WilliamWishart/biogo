//go:build tools
// +build tools

// This file ensures tool dependencies are tracked by Go modules.
package main

import (
	_ "github.com/cweill/gotests/..."
	_ "github.com/go-delve/delve/cmd/dlv"
	_ "github.com/haya14busa/goplay/cmd/goplay"
	_ "github.com/josharian/impl"
	_ "golang.org/x/tools/cmd/gopls"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
