package md2html

import (
	"context"
	"fmt"
	"io"
	_ "log/slog"
	"path/filepath"
	"strings"

	"github.com/aaronland/gocloud-blob/walk"
	"github.com/whosonfirst/go-blog/markdown"
	"github.com/whosonfirst/go-blog/parser"
	"github.com/whosonfirst/go-blog/render"
	"gocloud.dev/blob"
)

func RenderDirectory(ctx context.Context, html_opts *render.HTMLOptions, dir string) error {

	walk_func := func(ctx context.Context, obj *blob.ListObject) error {
		return RenderPath(ctx, html_opts, obj.Key)
	}

	return walk.WalkBucket(ctx, html_opts.SourceBucket, walk_func)
}

func RenderPath(ctx context.Context, html_opts *render.HTMLOptions, path string) error {

	select {

	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	fname := filepath.Base(path)

	if fname != html_opts.Input {
		return nil
	}

	r, err := html_opts.SourceBucket.NewReader(ctx, path, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new reader for %s, %w", path, err)
	}

	defer r.Close()

	parse_opts := parser.DefaultParseOptions()
	fm, body, err := parser.ParseReaderWithURI(ctx, parse_opts, r, path)

	if err != nil {
		return err
	}

	out_path := fm.Permalink

	if out_path == "" {
		root := filepath.Dir(path)
		out_path = filepath.Join(root, html_opts.Output)
	}

	if strings.HasSuffix(out_path, "/") {
		out_path = filepath.Join(out_path, html_opts.Output)
	}

	// START OF reconcile with RenderHTML in wof-md2idx

	doc, err := markdown.NewDocument(fm, body)
	html_r, err := render.RenderHTML(doc, html_opts)

	if err != nil {
		return err
	}

	defer html_r.Close()

	html_wr, err := html_opts.SourceBucket.NewWriter(ctx, out_path, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new writer for %s, %w", out_path, err)
	}

	_, err = io.Copy(html_wr, html_r)

	if err != nil {
		return fmt.Errorf("Failed to write %s, %w", out_path, err)
	}

	err = html_wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close %s after writing, %w", out_path, err)
	}

	// END OF reconcile with RenderHTML in wof-md2idx

	return nil
}

func Render(ctx context.Context, path string, html_opts *render.HTMLOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	switch html_opts.Mode {

	case "files":
		return RenderPath(ctx, html_opts, path)
	case "directory":
		return RenderDirectory(ctx, html_opts, path)
	default:
		return fmt.Errorf("Unknown or invalid mode")
	}
}
