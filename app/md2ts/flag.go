package md2ts

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var md_bucket_uri string
var url_prefix string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("md2ts")

	fs.StringVar(&md_bucket_uri, "markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")
	fs.StringVar(&url_prefix, "url-prefix", "", "An option prefix to append to all record URLs in the tinysearch index.")
	
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Generate a tinysearch compatible JSON file of blog posts used to generate a tinysearch index.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	return fs
}
