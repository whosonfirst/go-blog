package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/aaronland/gocloud-blob/walk"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-blog/jekyll"
	"github.com/whosonfirst/go-blog/parser"
	"github.com/whosonfirst/go-blog/render"
	"github.com/whosonfirst/go-blog/templates"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
)

func RenderDirectory(ctx context.Context, opts *render.FeedOptions, uri string) error {

	posts, err := GatherPosts(ctx, opts, uri)

	if err != nil {
		return fmt.Errorf("Failed to gather posts for %s, %w", uri, err)
	}

	if len(posts) == 0 {
		return nil
	}

	return RenderPosts(ctx, opts, uri, posts)
}

// THIS IS A BAD NAME - ALSO SHOULD BE SHARED CODE...
// (20180130/thisisaaronland)

func RenderPath(ctx context.Context, opts *render.FeedOptions, uri string) (*jekyll.FrontMatter, error) {

	select {

	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	r, err := opts.SourceBucket.NewReader(ctx, uri, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for reading, %w", uri, err)
	}

	defer r.Close()

	parse_opts := parser.DefaultParseOptions()
	parse_opts.Body = false

	fm, _, err := parser.ParseReaderWithURI(ctx, parse_opts, r, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s, %w", uri, err)
	}

	return fm, nil
}

func GatherPosts(ctx context.Context, opts *render.FeedOptions, uri string) ([]*jekyll.FrontMatter, error) {

	mu := new(sync.Mutex)

	lookup := make(map[string]*jekyll.FrontMatter)
	dates := make([]string, 0)

	walk_func := func(ctx context.Context, obj *blob.ListObject) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		walk_uri := obj.Key

		if filepath.Base(walk_uri) != opts.Input {
			return nil
		}

		fm, err := RenderPath(ctx, opts, walk_uri)

		if err != nil {
			return err
		}

		if fm == nil {
			return nil
		}

		mu.Lock()
		ymd := fm.Date.Format("20060102")
		dates = append(dates, ymd)
		lookup[ymd] = fm
		mu.Unlock()
		return nil
	}

	err := walk.WalkBucket(ctx, opts.SourceBucket, walk_func)

	if err != nil {
		return nil, fmt.Errorf("Failed to walk %s, %w", uri, err)
	}

	posts := make([]*jekyll.FrontMatter, 0)

	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	for _, ymd := range dates {
		posts = append(posts, lookup[ymd])

		if len(posts) == opts.Items {
			break
		}
	}

	return posts, nil
}

func RenderPosts(ctx context.Context, opts *render.FeedOptions, uri string, posts []*jekyll.FrontMatter) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	type Data struct {
		Posts     []*jekyll.FrontMatter
		BuildDate time.Time
	}

	now := time.Now()

	d := Data{
		Posts:     posts,
		BuildDate: now,
	}

	out_path := filepath.Join(uri, opts.Output)

	wr, err := opts.TargetBucket.NewWriter(ctx, out_path, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new writer for %s, %w", out_path, err)
	}

	t_name := fmt.Sprintf("feed_%s", opts.Format)
	t := opts.Templates.Lookup(t_name)

	if t == nil {
		return fmt.Errorf("Invalid or missing template '%s'", t_name)
	}

	err = t.Execute(wr, d)

	if err != nil {
		return fmt.Errorf("Failed to render template for %s, %w", out_path, err)
	}

	err = wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close %s after writing, %w", out_path, err)
	}

	return nil
}

func main() {

	var input = flag.String("input", "index.md", "What you expect the input Markdown file to be called")
	var output = flag.String("output", "", "The filename of your feed. If empty default to the value of -format + \".xml\"")

	var format = flag.String("format", "rss_20", "Valid options are: atom_10, rss_20")
	var items = flag.Int("items", 10, "The number of items to include in your feed")

	var feeds_bucket_uri = flag.String("feeds-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where feeds should be written to.")
	var md_bucket_uri = flag.String("markdown-bucket-uri", "", "A valid gocloud.dev/blob bucket URI where Markdown files should be read from.")

	var templates_uris multi.MultiString
	flag.Var(&templates_uris, "template-uri", "One or more valid gocloud.dev/blob bucket URIs where feed template files should be read from.")

	flag.Parse()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t, err := templates.LoadFeedTemplates(ctx, templates_uris...)

	if err != nil {
		log.Fatalf("Failed to load HTML templates, %v", err)
	}

	md_bucket, err := bucket.OpenBucket(ctx, *md_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open Markdown bucket, %v", err)
	}

	defer md_bucket.Close()

	feeds_bucket, err := bucket.OpenBucket(ctx, *feeds_bucket_uri)

	if err != nil {
		log.Fatalf("Failed to open HTML bucket, %v", err)
	}

	if *output == "" {
		*output = fmt.Sprintf("%s.xml", *format)
	}

	opts := render.DefaultFeedOptions()
	opts.Input = *input
	opts.Output = *output
	opts.Format = *format
	opts.Items = *items
	opts.Templates = t
	opts.SourceBucket = md_bucket
	opts.TargetBucket = feeds_bucket

	for _, uri := range flag.Args() {

		err := RenderDirectory(ctx, opts, uri)

		if err != nil {
			log.Fatalf("Failed to render %s, %v", uri, err)
		}
	}
}
