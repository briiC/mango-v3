package mango

import (
	"bytes"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func parseMarkdown(buf []byte) []byte {
	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)
	buf = bytes.ReplaceAll(buf, []byte("\t"), []byte("    "))
	buf = bytes.ReplaceAll(buf, []byte("\r"))
	return markdown.ToHTML(buf, parser, nil)
}
