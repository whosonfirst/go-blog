package md2html

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var mode string
var input string
var output string
var header string
var footer string

var html_bucket_uri string
var md_bucket_uri string

var templates_uris multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("md2html")

	fs.StringVar(&mode, "mode", "files", "Valid modes are: files, directory")
	fs.StringVar(&input, "input", "index.md", "What you expect the input Markdown file to be called")
	fs.StringVar(&output, "output", "index.html", "What you expect the output HTML file to be called")
	fs.StringVar(&header, "header", "", "The name of the (Go) template to use as a custom header")
	fs.StringVar(&footer, "footer", "", "The name of the (Go) template to use as a custom footer")

	fs.StringVar(&html_bucket_uri, "html-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where HTML files should be written to.")
	fs.StringVar(&md_bucket_uri, "markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")

	fs.Var(&templates_uris, "template-uri", "One or more valid gocloud.dev/blob bucket URIs where HTML template files should be read from.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Converts a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and converts them to HTML documents writing them to a target gocloud.dev/blob bucket URI.\n")
		fmt.Fprintf(os.Stderr, "Parse one or more whosonfirst/go-blog -style Markdown URIs and output FrontMatter, body text or both.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	return fs
}
