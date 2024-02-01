package md2feed

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var input string
var output string
var format string
var items int

var feeds_bucket_uri string
var md_bucket_uri string

var templates_uris multi.MultiString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("md2feed")

	fs.StringVar(&input, "input", "index.md", "What you expect the input Markdown file to be called")
	fs.StringVar(&output, "output", "", "The filename of your feed. If empty default to the value of -format + \".xml\"")

	fs.StringVar(&format, "format", "rss_20", "Valid options are: atom_10, rss_20")
	fs.IntVar(&items, "items", 10, "The number of items to include in your feed")

	fs.StringVar(&feeds_bucket_uri, "feeds-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where feeds should be written to.")
	fs.StringVar(&md_bucket_uri, "markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")

	fs.Var(&templates_uris, "template-uri", "One or more valid gocloud.dev/blob bucket URIs where feed template files should be read from.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Generate Atom 1.0 or RSS 2.0 syndication feeds from a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and writing the feeds to a target gocloud.dev/blob bucket URI.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	return fs
}
