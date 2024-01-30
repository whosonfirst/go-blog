// wof-md2html converts a directory of Markdown files in to HTML.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/aaronland/gocloud-blob/walk"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-blog"
	"github.com/whosonfirst/go-blog/parser"
	"github.com/whosonfirst/go-blog/render"
	"github.com/whosonfirst/go-blog/templates"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
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

func main() {

	var mode = flag.String("mode", "files", "Valid modes are: files, directory")
	var input = flag.String("input", "index.md", "What you expect the input Markdown file to be called")
	var output = flag.String("output", "index.html", "What you expect the output HTML file to be called")
	var header = flag.String("header", "", "The name of the (Go) template to use as a custom header")
	var footer = flag.String("footer", "", "The name of the (Go) template to use as a custom footer")

	var html_bucket_uri = flag.String("html-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where HTML files should be written to.")
	var md_bucket_uri = flag.String("markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")

	var templates_uris multi.MultiString
	flag.Var(&templates_uris, "template-uri", "One or more valid gocloud.dev/blob bucket URIs where HTML template files should be read from.")

	flag.Parse()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t, err := templates.LoadHTMLTemplates(ctx, templates_uris...)

	if err != nil {
		log.Fatalf("Failed to load HTML templates, %v", err)
	}

	md_bucket, err := bucket.OpenBucket(ctx, *md_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open Markdown bucket, %v", err)
	}

	defer md_bucket.Close()

	html_bucket, err := bucket.OpenBucket(ctx, *html_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open HTML bucket, %v", err)
	}

	defer html_bucket.Close()

	opts := render.DefaultHTMLOptions()
	opts.Mode = *mode
	opts.Input = *input
	opts.Output = *output
	opts.Header = *header
	opts.Footer = *footer
	opts.Templates = t
	opts.SourceBucket = md_bucket
	opts.TargetBucket = html_bucket

	for _, path := range flag.Args() {

		err := Render(ctx, path, opts)

		if err != nil {
			log.Fatalf("Failed to render %s, %v", path, err)
		}
	}
}
