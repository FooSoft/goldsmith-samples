package main

import (
	"flag"
	"log"

	"foosoft.net/projects/goldsmith"
	"foosoft.net/projects/goldsmith-components/devserver"
	"foosoft.net/projects/goldsmith-components/filters/condition"
	"foosoft.net/projects/goldsmith-components/plugins/frontmatter"
	"foosoft.net/projects/goldsmith-components/plugins/layout"
	"foosoft.net/projects/goldsmith-components/plugins/markdown"
	"foosoft.net/projects/goldsmith-components/plugins/minify"
)

type builder struct {
	dist bool
}

func (self *builder) Build(srcDir, dstDir, cacheDir string) {
	errs := goldsmith.
		Begin(srcDir).                        // read files from srcDir
		Chain(frontmatter.New()).             // extract frontmatter and store it as metadata
		Chain(markdown.New()).                // convert *.md files to *.html files
		Chain(layout.New()).                  // apply *.gohtml templates to *.html files
		FilterPush(condition.New(self.dist)). // push a dist-only conditional filter onto the stack
		Chain(minify.New()).                  // minify *.html, *.css, *.js, etc. files
		FilterPop().                          // pop off the last filter pushed onto the stack
		End(dstDir)                           // write files to dstDir

	for _, err := range errs {
		log.Print(err)
	}
}

func main() {
	port := flag.Int("port", 8080, "server port")
	dist := flag.Bool("dist", false, "final dist mode")
	flag.Parse()

	devserver.DevServe(&builder{*dist}, *port, "content", "build", "cache")
}
