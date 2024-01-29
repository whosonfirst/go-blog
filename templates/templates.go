package templates

import (
	"context"
	"fmt"
	html_template "html/template"
	"io"
	"path/filepath"
	text_template "text/template"

	"github.com/aaronland/gocloud-blob/walk"
	"gocloud.dev/blob"
)

// These should not be necessary but passing *blob.Bucket to template.ParseFS triggers "no matching patterns" errors
// even though the files (with the matching patterns) are in fact in the bucket...

func parseHTMLTemplates(ctx context.Context, t *html_template.Template, b *blob.Bucket, ext string) (*html_template.Template, error) {

	walk_func := func(ctx context.Context, obj *blob.ListObject) error {

		if filepath.Ext(obj.Key) != ext {
			return nil
		}

		r, err := b.NewReader(ctx, obj.Key, nil)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", obj.Key, err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", obj.Key, err)
		}

		t, err = t.Parse(string(body))

		if err != nil {
			return fmt.Errorf("Failed to parse template for %s, %w", obj.Key, err)
		}

		return nil
	}

	err := walk.WalkBucket(ctx, b, walk_func)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse FS for bucket, %w", err)
	}

	return t, nil
}

func parseTextTemplates(ctx context.Context, t *text_template.Template, b *blob.Bucket, ext string) (*text_template.Template, error) {

	walk_func := func(ctx context.Context, obj *blob.ListObject) error {

		if filepath.Ext(obj.Key) != ext {
			return nil
		}

		r, err := b.NewReader(ctx, obj.Key, nil)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", obj.Key, err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", obj.Key, err)
		}

		t, err = t.Parse(string(body))

		if err != nil {
			return fmt.Errorf("Failed to parse template for %s, %w", obj.Key, err)
		}

		return nil
	}

	err := walk.WalkBucket(ctx, b, walk_func)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse FS for bucket, %w", err)
	}

	return t, nil
}
