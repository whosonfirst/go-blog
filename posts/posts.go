package posts

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"path/filepath"

	"github.com/aaronland/gocloud-blob/walk"
	"github.com/whosonfirst/go-blog/jekyll"
	"github.com/whosonfirst/go-blog/markdown"
	"github.com/whosonfirst/go-blog/parser"
	"gocloud.dev/blob"
)

type Post struct {
	FrontMatter *jekyll.FrontMatter
	Body        *markdown.Body
}

func Iterate(ctx context.Context, source_bucket *blob.Bucket) iter.Seq2[*Post, error] {

	return func(yield func(*Post, error) bool) {

		walk_func := func(ctx context.Context, obj *blob.ListObject) error {

			select {
			case <-ctx.Done():
				return nil
			default:
				// pass
			}

			path := obj.Key

			logger := slog.Default()
			logger = logger.With("path", path)

			ext := filepath.Ext(path)

			if ext != ".md" {
				logger.Debug("Skip")
				return nil
			}

			r, err := source_bucket.NewReader(ctx, path, nil)

			if err != nil {
				logger.Error("Failed to create new reader", "error", err)
				yield(nil, err)
				return nil
			}

			defer r.Close()

			parse_opts := parser.DefaultParseOptions()
			fm, body, err := parser.ParseReaderWithURI(ctx, parse_opts, r, path)

			if err != nil {
				logger.Error("Failed to parse Markdown", "error", err)
				return nil
			}

			if fm == nil {
				logger.Error("File is missing front matter")
				return nil
			}

			p := &Post{
				FrontMatter: fm,
				Body:        body,
			}

			yield(p, nil)
			return nil
		}

		err := walk.WalkBucket(ctx, source_bucket, walk_func)

		if err != nil {
			yield(nil, fmt.Errorf("Failed to walk source bucket, %w", err))
		}
	}
}
