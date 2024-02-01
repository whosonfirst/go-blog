package mdparse

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var frontmatter bool
var body bool
var all bool

var md_bucket_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("mdparse")

	fs.BoolVar(&frontmatter, "frontmatter", false, "Dump (Jekyll) frontmatter")
	fs.BoolVar(&body, "body", false, "Dump (Markdown) body")
	fs.BoolVar(&all, "all", false, "Dump both frontmatter and body")

	fs.StringVar(&md_bucket_uri, "markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Parse one or more whosonfirst/go-blog -style Markdown URIs and output FrontMatter, body text or both.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	return fs
}
