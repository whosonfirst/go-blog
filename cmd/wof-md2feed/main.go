// wof-md2feed generate Atom 1.0 or RSS 2.0 syndication feeds from a collection of Markdown documents read from a source gocloud.dev/blob bucket URI
// and writing the feeds to a target gocloud.dev/blob bucket URI.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/whosonfirst/go-blog/app/md2feed"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := md2feed.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to run md2feed", "error", err)
		os.Exit(0)
	}
}
