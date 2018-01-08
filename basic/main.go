package main

import (
	"log"

	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-devserver"
	"github.com/FooSoft/goldsmith-plugins/plugins/frontmatter"
	"github.com/FooSoft/goldsmith-plugins/plugins/layout"
	"github.com/FooSoft/goldsmith-plugins/plugins/markdown"
)

type builder struct{}

func (b *builder) Build(srcDir, dstDir string) {
	errs := goldsmith.Begin(srcDir).
		Chain(frontmatter.New()).
		Chain(markdown.New()).
		Chain(layout.New("layouts/*.html")).
		End(dstDir)

	for _, err := range errs {
		log.Print(err)
	}
}

func main() {
	devserver.DevServe(new(builder), 8080, "src", "dst", "layouts")
}
