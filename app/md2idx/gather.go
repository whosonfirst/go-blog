package md2idx

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/aaronland/gocloud-blob/walk"
	"github.com/whosonfirst/go-blog/jekyll"
	"github.com/whosonfirst/go-blog/render"
	"gocloud.dev/blob"
)

func GatherPosts(ctx context.Context, html_opts *render.HTMLOptions, md_opts *MarkdownOptions, uri string) (map[string][]*jekyll.FrontMatter, error) {

	mu := new(sync.Mutex)

	lookup := make(map[string][]*jekyll.FrontMatter)

	walk_func := func(ctx context.Context, obj *blob.ListObject) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		walk_uri := obj.Key

		if filepath.Base(walk_uri) != html_opts.Input {
			return nil
		}

		fm, err := FrontMatterForPath(ctx, html_opts, walk_uri)

		if err != nil {
			return err
		}

		if fm == nil {
			return nil
		}

		var keys []string

		switch md_opts.Mode {
		case "authors":
			keys = fm.Authors
		case "ymd":
			keys = []string{
				fm.Date.Format("2006/01/02"),
				fm.Date.Format("2006/01"),
				fm.Date.Format("2006"),
			}
		case "tags":
			keys = fm.Tags
		case "landing":
			keys = []string{
				fm.Date.Format("20060102"),
			}
		default:
			return fmt.Errorf("Invalid or unsupported mode '%s'", md_opts.Mode)
		}

		mu.Lock()

		for _, k := range keys {

			posts, ok := lookup[k]

			if ok {
				posts = append(posts, fm)
				lookup[k] = posts
			} else {
				posts = []*jekyll.FrontMatter{fm}
				lookup[k] = posts
			}

		}

		mu.Unlock()
		return nil
	}

	err := walk.WalkBucket(ctx, html_opts.SourceBucket, walk_func)

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
