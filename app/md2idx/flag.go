package md2idx

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var input string
var output string
var header string
var footer string
var list string
var rollup string
var mode string

var html_bucket_uri string
var md_bucket_uri string

var per_page int

var html_templates_uris multi.MultiString
var md_templates_uris multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("md2html")

	fs.StringVar(&input, "input", "index.md", "What you expect the input Markdown file to be called")
	fs.StringVar(&output, "output", "index.html", "What you expect the output HTML file to be called")
	fs.StringVar(&header, "header", "", "The name of the (Go) template to use as a custom header")
	fs.StringVar(&footer, "footer", "", "The name of the (Go) template to use as a custom footer")
	fs.StringVar(&list, "list", "", "The name of the (Go) template to use as a custom list view")
	fs.StringVar(&rollup, "rollup", "", "The name of the (Go) template to use as a custom rollup view (for things like tags and authors)")
	fs.StringVar(&mode, "mode", "landing", "Valid modes are: authors, landing, tags, ymd.")

	fs.Var(&html_templates_uris, "html-template-uri", "One or more valid gocloud.dev/blob bucket URIs where HTML template files should be read from.")
	fs.Var(&md_templates_uris, "markdown-template-uri", "One or more valid gocloud.dev/blob bucket URIs where Markdown template files should be read from.")

	fs.StringVar(&html_bucket_uri, "html-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where HTML files should be written to.")
	fs.StringVar(&md_bucket_uri, "markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")

	fs.IntVar(&per_page, "per-page", 10, "The number of posts to include on a single page (the rest will be paginated on to page2, page3 and so on.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Generate paginated \"index\"-style list pages for a collection of blog posts. List styles include authors, tags, dates and reverse-chronological posts.\n")
		fmt.Fprintf(os.Stderr, "Parse one or more whosonfirst/go-blog -style Markdown URIs and output FrontMatter, body text or both.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	return fs
}
