package templates

import (
	"context"
	"fmt"
	html_template "html/template"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/whosonfirst/go-blog/uri"
)

func LoadHTMLTemplates(ctx context.Context, bucket_uris ...string) (*html_template.Template, error) {

	var fns = html_template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
		"prune_string": uri.PruneString,
	}

	t := html_template.New("html").Funcs(fns)

	for _, uri := range bucket_uris {

		b, err := bucket.OpenBucket(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to open %s for walking, %w", uri, err)
		}

		defer b.Close()

		t, err = parseHTMLTemplates(ctx, t, b, ".html")

		if err != nil {
			return nil, fmt.Errorf("Failed to parse templates for %s, %w", uri, err)
		}
	}

	return t, nil
}
