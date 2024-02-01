package md2idx

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

	logger = logger.With("cmd", "md2idx")
	slog.SetDefault(logger)

	md_bucket, err := bucket.OpenBucket(ctx, md_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open Markdown bucket, %w", err)
	}

	defer md_bucket.Close()

	html_bucket, err := bucket.OpenBucket(ctx, html_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open HTML bucket, %w", err)
	}

	html_t, err := templates.LoadHTMLTemplates(ctx, html_templates_uris...)

	if err != nil {
		return fmt.Errorf("Failed to load HTML templates, %w", err)
	}

	md_t, err := templates.LoadMarkdownTemplates(ctx, md_templates_uris...)

	if err != nil {
		return fmt.Errorf("Failed to load Markdown templates, %w", err)
	}

	html_opts := render.DefaultHTMLOptions()
	html_opts.Input = input
	html_opts.Output = output
	html_opts.Header = header
	html_opts.Footer = footer
	html_opts.Templates = html_t
	html_opts.SourceBucket = md_bucket
	html_opts.TargetBucket = html_bucket
	html_opts.PerPage = per_page

	md_opts := &MarkdownOptions{
		MarkdownTemplates: md_t,
		List:              list,
		Rollup:            rollup,
		Mode:              mode,
	}

	for _, uri := range fs.Args() {

		err := RenderDirectory(ctx, html_opts, md_opts, uri)

		if err != nil {
			return fmt.Errorf("Failed to render %s, %w", uri, err)
		}
	}

	return nil
}
