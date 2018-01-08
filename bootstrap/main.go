/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"flag"
	"log"

	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-components/devserver"
	"github.com/FooSoft/goldsmith-components/filters/condition"
	"github.com/FooSoft/goldsmith-components/plugins/abs"
	"github.com/FooSoft/goldsmith-components/plugins/breadcrumbs"
	"github.com/FooSoft/goldsmith-components/plugins/collection"
	"github.com/FooSoft/goldsmith-components/plugins/dom"
	"github.com/FooSoft/goldsmith-components/plugins/frontmatter"
	"github.com/FooSoft/goldsmith-components/plugins/index"
	"github.com/FooSoft/goldsmith-components/plugins/layout"
	"github.com/FooSoft/goldsmith-components/plugins/markdown"
	"github.com/FooSoft/goldsmith-components/plugins/paginate"
	"github.com/FooSoft/goldsmith-components/plugins/syntax"
	"github.com/FooSoft/goldsmith-components/plugins/tags"
	"github.com/FooSoft/goldsmith-components/plugins/thumbnail"
	"github.com/PuerkitoBio/goquery"
)

func fixup(doc *goquery.Document) error {
	doc.Find("table").AddClass("table")
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

type builder struct {
	root string
}

func (b *builder) Build(srcDir, dstDir string) {
	tagMeta := map[string]interface{}{
		"Area":        "tags",
		"CrumbParent": "tags",
		"Layout":      "tag",
	}

	indexMeta := map[string]interface{}{
		"Layout": "index",
	}

	errs := goldsmith.Begin(srcDir).
		Chain(frontmatter.New()).
		Chain(markdown.New().SummaryKey("Summary")).
		Chain(collection.New()).
		Chain(paginate.New("Groups.Blog").InheritKeys("Layout").ItemsPerPage(3)).
		Chain(index.New(indexMeta)).
		Chain(tags.New().IndexMeta(tagMeta)).
		Chain(breadcrumbs.New()).
		Chain(layout.New("layouts/*.html")).
		Chain(syntax.New().Placement(syntax.PlaceInline)).
		Chain(dom.New(fixup)).
		Chain(thumbnail.New()).
		FilterPush(condition.New(len(b.root) > 0)).
		Chain(abs.New().BaseUrl(b.root)).
		FilterPop().
		End(dstDir)

	for _, err := range errs {
		log.Print(err)
	}
}

func main() {
	var (
		root = flag.String("root", "", "root directory")
		port = flag.Int("port", 8080, "server port")
	)
	flag.Parse()

	devserver.DevServe(&builder{*root}, *port, "src", "dst", "layouts")
}
