// wof-md2html converts a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and converts
// them to HTML documents writing them to a target gocloud.dev/blob bucket URI.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/whosonfirst/go-blog/app/md2html"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := md2html.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to run md2html", "error", err)
		os.Exit(0)
	}
}
