package md2feed

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

	t, err := templates.LoadFeedTemplates(ctx, templates_uris...)

	if err != nil {
		return fmt.Errorf("Failed to load HTML templates, %v", err)
	}

	md_bucket, err := bucket.OpenBucket(ctx, md_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open Markdown bucket, %v", err)
	}

	defer md_bucket.Close()

	feeds_bucket, err := bucket.OpenBucket(ctx, feeds_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open HTML bucket, %v", err)
	}

	if output == "" {
		output = fmt.Sprintf("%s.xml", format)
	}

	opts := render.DefaultFeedOptions()
	opts.Input = input
	opts.Output = output
	opts.Format = format
	opts.Items = items
	opts.Templates = t
	opts.SourceBucket = md_bucket
	opts.TargetBucket = feeds_bucket

	for _, uri := range flag.Args() {

		err := RenderDirectory(ctx, opts, uri)

		if err != nil {
			return fmt.Errorf("Failed to render %s, %v", uri, err)
		}
	}

	return nil
}
