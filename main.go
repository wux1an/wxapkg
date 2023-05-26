package main

import (
	"fmt"
	"github.com/wux1an/wxapkg/cmd"
)

var version = "v0.0.1"
var commit = "b7122099"

func main() {
	cmd.RootCmd.Version = fmt.Sprintf("%s(%s)", version, commit)
	cmd.Execute()
}
