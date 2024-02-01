// wof-md2idx will generate paginated "index"-style list pages for a collection of blog posts. List styles include authors,
// tags, dates and reverse-chronological posts.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/whosonfirst/go-blog/app/md2idx"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := md2idx.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to run md2idx", "error", err)
		os.Exit(0)
	}
}
