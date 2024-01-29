package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/whosonfirst/go-blog/parser"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
)

func main() {

	var frontmatter = flag.Bool("frontmatter", false, "Dump (Jekyll) frontmatter")
	var body = flag.Bool("body", false, "Dump (Markdown) body")
	var all = flag.Bool("all", false, "Dump both frontmatter and body")

	var md_bucket_uri = flag.String("markdown-bucket-uri", "cwd://", "...")

	flag.Parse()

	ctx := context.Background()

	md_bucket, err := bucket.OpenBucket(ctx, *md_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open Markdown bucket, %v", err)
	}

	defer md_bucket.Close()

	if *all {
		*frontmatter = true
		*body = true
	}

	opts := parser.DefaultParseOptions()
	opts.FrontMatter = *frontmatter
	opts.Body = *body

	for _, uri := range flag.Args() {

		r, err := md_bucket.NewReader(ctx, uri, nil)

		if err != nil {
			log.Fatalf("Failed to create new reader for %s, %v", uri, err)
		}

		defer r.Close()

		fm, b, err := parser.ParseReaderWithURI(ctx, opts, r, uri)

		if err != nil {
			log.Fatalf("Failed to parse %s, %v", uri, err)
		}

		if *frontmatter {
			fmt.Println(fm.String())
		}

		if *body {
			fmt.Println(b.String())
		}

	}
}
