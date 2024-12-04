package md2ts

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/aaronland/gocloud-blob/walk"
	"github.com/whosonfirst/go-blog/jekyll"
	"github.com/whosonfirst/go-blog/parser"
	"gocloud.dev/blob"
)

func GatherPosts(ctx context.Context, source_bucket *blob.Bucket) (map[string][]*jekyll.FrontMatter, error) {

	mu := new(sync.Mutex)

	lookup := make(map[string][]*jekyll.FrontMatter)

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
			return fmt.Errorf("Failed to create new reader for %s, %w", path, err)
		}

		defer r.Close()

		parse_opts := parser.DefaultParseOptions()
		fm, _, err := parser.ParseReaderWithURI(ctx, parse_opts, r, path)

		if err != nil {
			logger.Error("Failed to parse Markdown", "error", err)
			return nil
		}

		if fm == nil {
			logger.Error("File is missing front matter")
			return nil
		}

		mu.Lock()

		lookup[path] = []*jekyll.FrontMatter{
			fm,
		}

		mu.Unlock()
		return nil
	}

	err := walk.WalkBucket(ctx, source_bucket, walk_func)

	if err != nil {
		return nil, fmt.Errorf("Failed to walk source bucket, %w", err)
	}

	// ensure that everything is sorted by date (reverse chronological)

	for k, unsorted := range lookup {

		count := len(unsorted)

		by_date := make(map[string]*jekyll.FrontMatter)
		dates := make([]string, count)

		for idx, post := range unsorted {

			dt := post.Date.Format(time.RFC3339)

			by_date[dt] = post
			dates[idx] = dt
		}

		sort.Sort(sort.Reverse(sort.StringSlice(dates)))

		sorted := make([]*jekyll.FrontMatter, count)

		for idx, dt := range dates {
			sorted[idx] = by_date[dt]
		}

		lookup[k] = sorted
	}

	return lookup, nil
}
