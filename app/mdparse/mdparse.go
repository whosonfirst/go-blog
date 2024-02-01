package mdparse

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-blog/parser"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	flagset.Parse(fs)

	md_bucket, err := bucket.OpenBucket(ctx, md_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open Markdown bucket, %w", err)
	}

	defer md_bucket.Close()

	if all {
		frontmatter = true
		body = true
	}

	opts := parser.DefaultParseOptions()
	opts.FrontMatter = frontmatter
	opts.Body = body

	for _, uri := range flag.Args() {

		r, err := md_bucket.NewReader(ctx, uri, nil)

		if err != nil {
			return fmt.Errorf("Failed to create new reader for %s, %w", uri, err)
		}

		defer r.Close()

		fm, b, err := parser.ParseReaderWithURI(ctx, opts, r, uri)

		if err != nil {
			return fmt.Errorf("Failed to parse %s, %w", uri, err)
		}

		if frontmatter {
			fmt.Println(fm.String())
		}

		if body {
			fmt.Println(b.String())
		}
	}

	return nil
}
