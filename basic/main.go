package main

import (
	"flag"
	"log"

	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-components/devserver"
	"github.com/FooSoft/goldsmith-components/plugins/frontmatter"
	"github.com/FooSoft/goldsmith-components/plugins/layout"
	"github.com/FooSoft/goldsmith-components/plugins/markdown"
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
	port := flag.Int("port", 8080, "server port")
	flag.Parse()
	devserver.DevServe(new(builder), *port, "src", "dst", "layouts")
}
