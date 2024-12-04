package md2ts

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-blog/posts"
	"github.com/whosonfirst/go-blog/tinysearch"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	logger := slog.Default()
	logger = logger.With("cmd", "md2idx")

	md_bucket, err := bucket.OpenBucket(ctx, md_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open Markdown bucket, %w", err)
	}

	defer md_bucket.Close()

	index := make([]*tinysearch.Record, 0)

	for p, err := range posts.Iterate(ctx, md_bucket) {

		if err != nil {
			slog.Error("Failed to iterate", "error", err)
			break
		}

		r := &tinysearch.Record{
			Title: p.FrontMatter.Title,
			URL:   p.FrontMatter.Permalink,
			Body:  p.Body.String(),
		}

		index = append(index, r)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(index)

	if err != nil {
		return fmt.Errorf("Failed to encode index, %w", err)
	}

	return nil
}
