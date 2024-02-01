package md2html

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-blog/render"
	"github.com/whosonfirst/go-blog/templates"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	flagset.Parse(fs)

	logger = logger.With("cmd", "md2html")
	slog.SetDefault(logger)

	t, err := templates.LoadHTMLTemplates(ctx, templates_uris...)

	if err != nil {
		return fmt.Errorf("Failed to load HTML templates, %v", err)
	}

	md_bucket, err := bucket.OpenBucket(ctx, md_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open Markdown bucket, %v", err)
	}

	defer md_bucket.Close()

	html_bucket, err := bucket.OpenBucket(ctx, html_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open HTML bucket, %v", err)
	}

	defer html_bucket.Close()

	opts := render.DefaultHTMLOptions()
	opts.Mode = mode
	opts.Input = input
	opts.Output = output
	opts.Header = header
	opts.Footer = footer
	opts.Templates = t
	opts.SourceBucket = md_bucket
	opts.TargetBucket = html_bucket

	for _, path := range fs.Args() {

		err := Render(ctx, path, opts)

		if err != nil {
			return fmt.Errorf("Failed to render %s, %w", path, err)
		}
	}

	return nil
}
