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

func (b *builder) Build(contentDir, buildDir, cacheDir string) {
	errs := goldsmith.Begin(contentDir).
		Cache(cacheDir).
		Chain(frontmatter.New()).
		Chain(markdown.New()).
		Chain(layout.New()).
		End(buildDir)

	for _, err := range errs {
		log.Print(err)
	}
}

func main() {
	port := flag.Int("port", 8080, "server port")
	flag.Parse()

	devserver.DevServe(new(builder), *port, "content", "build", "cache")
}
