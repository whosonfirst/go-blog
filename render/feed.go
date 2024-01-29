package render

import (
	"text/template"

	"gocloud.dev/blob"
)

type FeedOptions struct {
	Format       string
	Input        string
	Output       string
	Items        int
	Templates    *template.Template
	SourceBucket *blob.Bucket
	TargetBucket *blob.Bucket
}

func DefaultFeedOptions() *FeedOptions {

	opts := FeedOptions{
		Input:     "index.md",
		Format:    "rss_20",
		Output:    "rss_20.xml",
		Items:     10,
		Templates: nil,
	}

	return &opts
}
