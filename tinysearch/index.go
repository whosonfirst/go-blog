package tinysearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"net/url"
	
	"github.com/whosonfirst/go-blog/posts"
)

// IndexPostsOptions defines configuration details for generating a collection of `Records` derived from a collection of posts and write them to an `io.Writer` instance.
type IndexPostsOptions struct {
	// Iterator is an `iter.Seq2` instance to yield `posts.Post` instances to index.
	Iterator iter.Seq2[*posts.Post, error]
	// Writer is where the final set of `Record` documents should be written.
	Writer   io.Writer
	// URLPrefix is an optional prefix to assign the `Record.URL` property for each post.
	URLPrefix string
}

// IndexPosts will generate a collection of `Records` derived from a collection of posts and write them to an `io.Writer` instance.
func IndexPosts(ctx context.Context, opts *IndexPostsOptions) error {

	index := make([]*Record, 0)

	for p, err := range opts.Iterator {

		if err != nil {
			slog.Error("Failed to iterate", "error", err)
			break
		}

		r_url := p.FrontMatter.Permalink
		
		if opts.URLPrefix != "" {
			
			v, err := url.JoinPath(opts.URLPrefix, r_url)

			if err != nil {
				slog.Error("Failed to assign URL prefix", "url", r_url, "error", err)
				break
			}

			r_url = v
		}
		
		r := &Record{
			Title: p.FrontMatter.Title,
			URL:   r_url,
			Body:  p.Body.String(),	// TBD: strip MD...
		}

		index = append(index, r)
	}

	enc := json.NewEncoder(opts.Writer)
	err := enc.Encode(index)

	if err != nil {
		return fmt.Errorf("Failed to encode index, %w", err)
	}

	return nil
}
