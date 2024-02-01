// wof-mdparse parses one or more whosonfirst/go-blog -style Markdown URIs and output FrontMatter, body text or both.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/whosonfirst/go-blog/app/mdparse"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := mdparse.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to run mdparse", "error", err)
		os.Exit(0)
	}
}
