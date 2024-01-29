package templates

import (
	"context"
	"fmt"
	text_template "text/template"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/whosonfirst/go-blog/uri"
)

func LoadMarkdownTemplates(ctx context.Context, bucket_uris ...string) (*text_template.Template, error) {

	var fns = text_template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
		"prune_string": uri.PruneString,
	}

	t := text_template.New("markdown").Funcs(fns)

	for _, uri := range bucket_uris {

		b, err := bucket.OpenBucket(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to open %s for walking, %w", uri, err)
		}

		defer b.Close()

		t, err = parseTextTemplates(ctx, t, b, ".md")

		if err != nil {
			return nil, fmt.Errorf("Failed to parse templates for %s, %w", uri, err)
		}

	}

	return t, nil
}
