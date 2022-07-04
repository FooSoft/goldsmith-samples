package main

import (
	"flag"
	"log"

	"foosoft.net/projects/goldsmith"
	"foosoft.net/projects/goldsmith-components/devserver"
	"foosoft.net/projects/goldsmith-components/plugins/absolute"
	"foosoft.net/projects/goldsmith-components/plugins/breadcrumbs"
	"foosoft.net/projects/goldsmith-components/plugins/collection"
	"foosoft.net/projects/goldsmith-components/plugins/document"
	"foosoft.net/projects/goldsmith-components/plugins/frontmatter"
	"foosoft.net/projects/goldsmith-components/plugins/index"
	"foosoft.net/projects/goldsmith-components/plugins/layout"
	"foosoft.net/projects/goldsmith-components/plugins/markdown"
	"foosoft.net/projects/goldsmith-components/plugins/summary"
	"foosoft.net/projects/goldsmith-components/plugins/syntax"
	"foosoft.net/projects/goldsmith-components/plugins/tags"
	"foosoft.net/projects/goldsmith-components/plugins/thumbnail"
	"github.com/PuerkitoBio/goquery"
)

func fixup(file *goldsmith.File, doc *goquery.Document) error {
	doc.Find("table").AddClass("table").Find("thead").AddClass("thead-light")
	doc.Find("blockquote").AddClass("blockquote")
	doc.Find("img[src*='thumb']").Each(func(i int, s *goquery.Selection) {
		thumbLink := s.ParentFiltered("a")
		thumbLink.AddClass("img-thumbnail", "img-thumbnail-inline")
		thumbLink.SetAttr("data-title", s.AttrOr("alt", ""))
		thumbLink.SetAttr("data-toggle", "lightbox")
		thumbLink.SetAttr("data-gallery", "gallery")
	})

	return nil
}

type builder struct{}

func (self *builder) Build(contentDir, buildDir, cacheDir string) {
	tagMeta := map[string]interface{}{
		"Area":        "tags",
		"CrumbParent": "tags",
		"Layout":      "tag",
	}

	indexMeta := map[string]interface{}{
		"Layout": "index",
	}

	errs := goldsmith.Begin(contentDir).
		Cache(cacheDir).
		Chain(frontmatter.New()).
		Chain(markdown.New()).
		Chain(summary.New()).
		Chain(collection.New()).
		Chain(index.New(indexMeta)).
		Chain(tags.New().IndexMeta(tagMeta)).
		Chain(breadcrumbs.New()).
		Chain(layout.New()).
		Chain(syntax.New().Placement(syntax.PlaceInline)).
		Chain(document.New(fixup)).
		Chain(thumbnail.New()).
		Chain(absolute.New()).
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
