// +build ignore

package main

import (
	"log"

	"github.com/jayemen/devcertsrv/assets"
	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(assets.Assets, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "assets/assets_embedded.go",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
