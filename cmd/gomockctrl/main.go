package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/thejan2009/gomockctrl"
)

func main() {
	singlechecker.Main(gomockctrl.Analyzer)
}
