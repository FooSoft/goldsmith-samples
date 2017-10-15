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
	"log"

	"github.com/FooSoft/goldsmith"
	"github.com/FooSoft/goldsmith-devserver"
	"github.com/FooSoft/goldsmith-plugins/breadcrumbs"
	"github.com/FooSoft/goldsmith-plugins/frontmatter"
	"github.com/FooSoft/goldsmith-plugins/index"
	"github.com/FooSoft/goldsmith-plugins/layout"
	"github.com/FooSoft/goldsmith-plugins/markdown"
	"github.com/FooSoft/goldsmith-plugins/syntax"
	"github.com/FooSoft/goldsmith-plugins/tags"
	"github.com/FooSoft/goldsmith-plugins/thumbnail"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
)

type fixup struct{}

func (*fixup) Name() string {
	return "fixup"
}

func (*fixup) Initialize(ctx goldsmith.Context) ([]string, error) {
	return []string{"**/*.html", "**/*.htm"}, nil
}

func (*fixup) Process(ctx goldsmith.Context, f goldsmith.File) error {
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return err
	}

	doc.Find("table").AddClass("table")
	doc.Find("blockquote").AddClass("blockquote")
	doc.Find("img[src*='thumb']").Each(func(i int, s *goquery.Selection) {
		thumbLink := s.ParentFiltered("a")
		thumbLink.AddClass("img-thumbnail", "img-thumbnail-inline")
		thumbLink.SetAttr("data-title", s.AttrOr("alt", ""))
		thumbLink.SetAttr("data-toggle", "lightbox")
		thumbLink.SetAttr("data-gallery", "gallery")
	})

	html, err := doc.Html()
	if err != nil {
		return err
	}

	nf := goldsmith.NewFileFromData(f.Path(), []byte(html))
	nf.CopyValues(f)
	ctx.DispatchFile(nf)
	return nil
}

type builder struct{}

func (b *builder) Build(srcDir, dstDir string) {
	tagMeta := map[string]interface{}{
		"Area":        "tags",
		"CrumbParent": "tags",
		"Layout":      "tag",
	}

	indexMeta := map[string]interface{}{
		"Layout": "index",
	}

	markdownFlags := blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK |
		blackfriday.EXTENSION_AUTO_HEADER_IDS

	errs := goldsmith.Begin(srcDir).
		Chain(frontmatter.New()).
		Chain(markdown.New().MarkdownFlags(markdownFlags)).
		Chain(index.New(indexMeta)).
		Chain(tags.New().IndexMeta(tagMeta)).
		Chain(breadcrumbs.New()).
		Chain(layout.New("layouts/*.html")).
		Chain(syntax.New().Placement(syntax.PlaceInline)).
		Chain(new(fixup)).
		Chain(thumbnail.New()).
		End(dstDir)

	for _, err := range errs {
		log.Print(err)
	}
}

func main() {
	devserver.DevServe(new(builder), 8080, "src", "dst", "layouts")
}
