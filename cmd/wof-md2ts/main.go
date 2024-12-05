package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-blog/app/md2ts"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()
	err := md2ts.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run md2ts, %v", err)
	}
}
