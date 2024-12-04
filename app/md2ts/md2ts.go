package md2ts

/*

$> go run cmd/wof-md2ts/main.go \
	-markdown-bucket-uri file:///usr/local/sfomuseum/www-sfomuseum-weblog/www/blog \
	-url-prefix http://millsfield.sfomuseum.org \
	> work/index.json

$> cd work
$> tinysearch index.json
$> fileserver -root ./wasm_output

*/

import (
	"context"
	"flag"
	"fmt"
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

	md_bucket, err := bucket.OpenBucket(ctx, md_bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open Markdown bucket, %w", err)
	}

	defer md_bucket.Close()

	posts_iter := posts.Iterate(ctx, md_bucket)
	wr := os.Stdout

	index_opts := &tinysearch.IndexPostsOptions{
		Iterator: posts_iter,
		Writer:   wr,
		URLPrefix: url_prefix,
	}

	err = tinysearch.IndexPosts(ctx, index_opts)

	if err != nil {
		return fmt.Errorf("Failed to index posts, %w", err)
	}

	return nil
}
