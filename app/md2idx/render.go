package md2idx

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/whosonfirst/go-blog/jekyll"
	"github.com/whosonfirst/go-blog/markdown"
	"github.com/whosonfirst/go-blog/parser"
	"github.com/whosonfirst/go-blog/render"
	wof_uri "github.com/whosonfirst/go-blog/uri"
)

var re_ymd *regexp.Regexp

var default_index_list string
var default_index_rollup string

func init() {

	re_ymd = regexp.MustCompile(".*(\\d{4})(?:/(\\d{2}))?(?:/(\\d{2}))?$")

	default_index_rollup = `{{ range $w := .Rollup}}
* [ {{ $w }} ]( {{ prune_string $w }} )
{{ end }}`

	default_index_list = `{{ range $fm := .Posts }}
### [{{ $fm.Title }}]({{ $fm.Permalink }}) 

> {{ $fm.Excerpt }}

{{$lena := len $fm.Authors }}
{{$lent := len $fm.Tags }}
<small class="this-is">This is a blog post by
    {{ range $ia, $a := $fm.Authors }}{{ if gt $lena 1 }}{{if eq $ia 0}}{{else if eq (plus1 $ia) $lena}} and {{else}}, {{end}}{{ end }}[{{ $a }}](/blog/authors/{{ prune_string $a  }}/){{ end }}.
    {{ if $fm.Date }}It was published on <span class="pubdate"><a href="/blog/{{ $fm.Date.Year }}/{{ $fm.Date.Format "01" }}/">{{ $fm.Date.Format "January" }}</a> <a href="/blog/{{ $fm.Date.Year }}/{{ $fm.Date.Format "01" }}/{{ $fm.Date.Format "02" }}/">{{ $fm.Date.Format "02"}}</a>, <a href="/blog/{{ $fm.Date.Year }}/">{{ $fm.Date.Format "2006" }}</a></span>{{ if gt $lent 0 }} and tagged {{ range $it, $t := $fm.Tags }}{{ if gt $lent 1 }}{{if eq $it 0}}{{else if eq (plus1 $it) $lent}} and {{else}}, {{end}}{{ end }}[{{ $t }}](/blog/tags/{{ prune_string $t  }}/){{ end }}{{ end}}.
    {{ else }}
    It was tagged {{ range $it, $t := $fm.Tags }}{{ if gt $lent 1 }}{{if eq $it 0}}{{else if eq (plus1 $it) $lent}} and {{else}}, {{end}}{{ end }}[{{ $t }}](/blog/tags/{{ prune_string $t  }}/){{ end }}.
    {{ end }}
</small>
{{ end }}`
}

type MarkdownOptions struct {
	MarkdownTemplates *template.Template
	List              string
	Rollup            string
	Mode              string
}

func RenderDirectory(ctx context.Context, html_opts *render.HTMLOptions, md_opts *MarkdownOptions, uri string) error {

	lookup, err := GatherPosts(ctx, html_opts, md_opts, uri)

	if err != nil {
		return fmt.Errorf("Failed to gather posts in %s, %w", uri, err)
	}

	keys := make([]string, 0)

	for k, _ := range lookup {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	if md_opts.Mode == "landing" {

		posts := make([]*jekyll.FrontMatter, 0)

		for _, k := range keys {

			for _, p := range lookup[k] {
				posts = append(posts, p)
			}
		}

		if len(posts) == 0 {
			return nil
		}

		title := "" // where is date...

		return RenderPosts(ctx, html_opts, md_opts, uri, title, posts)
	}

	var root string

	switch md_opts.Mode {
	case "authors", "tags":
		root = filepath.Join(uri, md_opts.Mode)
	case "ymd":
		root = uri
	default:
		return fmt.Errorf("Invalid or unsupported mode '%s'", md_opts.Mode)
	}

	for _, raw := range keys {

		var clean string

		if md_opts.Mode == "ymd" {
			clean = raw
		} else {

			c, err := wof_uri.PruneString(raw)

			if err != nil {
				return fmt.Errorf("Failed to clean string '%s', %w", raw, err)
			}

			clean = c
		}

		if clean == "" {
			continue
		}

		// html_opts.Title = raw

		k_dir := filepath.Join(root, clean)

		title := raw
		posts := lookup[raw]

		err = RenderPosts(ctx, html_opts, md_opts, k_dir, title, posts)

		if err != nil {
			return err
		}
	}

	switch md_opts.Mode {
	case "ymd":
		return nil
	default:
		return RenderRollup(ctx, root, keys, html_opts, md_opts)
	}
}

// see notes below about passing a struct for post details

func RenderPosts(ctx context.Context, html_opts *render.HTMLOptions, md_opts *MarkdownOptions, root string, title string, posts []*jekyll.FrontMatter) error {

	per_page := html_opts.PerPage
	count_posts := len(posts)

	pages := int(math.Ceil(float64(count_posts) / float64(per_page)))

	for i := 0; i < pages; i++ {

		from := i * per_page
		to := from + (per_page - 1)

		if to > count_posts {
			to = count_posts
		}

		if to == 0 {
			to = 1
		}

		pg := i + 1

		slice_posts := posts[from:to]

		pagination := &render.Pagination{
			Total: count_posts,
			Pages: pages,
			Page:  pg,
		}

		if pg > 1 {
			pagination.Previous = fmt.Sprintf("page%d.html", (pg - 1))
		}

		if pg < pages {
			pagination.Next = fmt.Sprintf("page%d.html", (pg + 1))
		}

		html_opts.Output = fmt.Sprintf("page%d.html", pg)
		html_opts.Pagination = pagination

		// log.Println(pg, root, html_opts.Output)

		err := renderPosts(ctx, root, title, slice_posts, html_opts, md_opts)

		if err != nil {
			return err
		}

		if i == 0 {

			// log.Println(pg, root, html_opts.Output)

			html_opts.Output = "index.html"
			err := renderPosts(ctx, root, title, slice_posts, html_opts, md_opts)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func renderPosts(ctx context.Context, root string, title string, posts []*jekyll.FrontMatter, html_opts *render.HTMLOptions, md_opts *MarkdownOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	var t *template.Template

	if md_opts.List != "" {
		t = md_opts.MarkdownTemplates.Lookup(md_opts.List)
	}

	if t == nil {

		func_map := template.FuncMap{
			"prune_string": wof_uri.PruneString,
			"plus1": func(x int) int {
				return x + 1
			},
		}

		tm, err := template.New("list").Funcs(func_map).Parse(default_index_list)

		if err != nil {
			return err
		}

		t = tm
	}

	// maybe just pass this to RenderPosts?
	// (20190409/thisisaaronland)

	type Data struct {
		Mode       string
		Title      string
		Posts      []*jekyll.FrontMatter
		Pagination *render.Pagination
	}

	d := Data{
		Mode:       md_opts.Mode,
		Title:      title,
		Posts:      posts,
		Pagination: html_opts.Pagination,
	}

	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	err := t.Execute(wr, d)

	if err != nil {
		return err
	}

	wr.Flush()

	r := bytes.NewReader(b.Bytes())
	fh := ioutil.NopCloser(r)

	parse_opts := parser.DefaultParseOptions()
	fm, buf, err := parser.Parse(fh, parse_opts)

	if err != nil {
		return fmt.Errorf("Failedto parse MD document, because %s\n", err)
	}

	if re_ymd.MatchString(root) {

		matches := re_ymd.FindStringSubmatch(root)

		str_yyyy := matches[1]
		str_mm := matches[2]
		str_dd := matches[3]

		parse_string := make([]string, 0)
		ymd_string := make([]string, 0)

		if str_yyyy != "" {
			parse_string = append(parse_string, "2006")
			ymd_string = append(ymd_string, str_yyyy)
		}

		if str_mm != "" {
			parse_string = append(parse_string, "01")
			ymd_string = append(ymd_string, str_mm)
		}

		if str_dd != "" {
			parse_string = append(parse_string, "02")
			ymd_string = append(ymd_string, str_dd)
		}

		// Y U SO WEIRD GO...

		dt, err := time.Parse(strings.Join(parse_string, "-"), strings.Join(ymd_string, "-"))

		if err == nil {
			fm.Date = &dt
		}
	}

	err = RenderHTML(ctx, root, html_opts, fm, buf)

	if err != nil {
		return err
	}

	return nil
}

func RenderRollup(ctx context.Context, root string, rollup []string, html_opts *render.HTMLOptions, md_opts *MarkdownOptions) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	var t *template.Template

	if md_opts.Rollup != "" {
		t = md_opts.MarkdownTemplates.Lookup(md_opts.Rollup)
	}

	if t == nil {

		func_map := template.FuncMap{
			"prune_string": wof_uri.PruneString,
			"plus1": func(x int) int {
				return x + 1
			},
		}

		tm, err := template.New("rollup").Funcs(func_map).Parse(default_index_rollup)

		if err != nil {
			return err
		}

		t = tm
	}

	sort.Sort(sort.StringSlice(rollup))

	type Data struct {
		Mode   string
		Rollup []string
	}

	d := Data{
		Mode:   md_opts.Mode,
		Rollup: rollup,
	}

	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	err := t.Execute(wr, d)

	if err != nil {
		return err
	}

	wr.Flush()

	r := bytes.NewReader(b.Bytes())
	fh := ioutil.NopCloser(r)

	parse_opts := parser.DefaultParseOptions()
	fm, buf, err := parser.Parse(fh, parse_opts)

	if err != nil {
		return fmt.Errorf("Failed to parse MD document, because %w\n", err)
	}

	html_opts.Output = "index.html"

	html_opts.Pagination = &render.Pagination{
		Total: len(rollup),
		Pages: 1,
		Page:  1,
	}

	err = RenderHTML(ctx, root, html_opts, fm, buf)

	if err != nil {
		return err
	}

	return nil
}

func RenderHTML(ctx context.Context, root string, html_opts *render.HTMLOptions, fm *jekyll.FrontMatter, body *markdown.Body) error {

	out_path := filepath.Join(root, html_opts.Output)

	doc, err := markdown.NewDocument(fm, body)

	if err != nil {
		return fmt.Errorf("Failed to create MD document, %w", err)
	}

	html_r, err := render.RenderHTML(doc, html_opts)

	if err != nil {
		return fmt.Errorf("Failed to render HTML, %w", err)
	}

	defer html_r.Close()

	html_wr, err := html_opts.TargetBucket.NewWriter(ctx, out_path, nil)

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

	return nil
}

func FrontMatterForPath(ctx context.Context, html_opts *render.HTMLOptions, uri string) (*jekyll.FrontMatter, error) {

	select {

	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	r, err := html_opts.SourceBucket.NewReader(ctx, uri, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for reading, %w", uri, err)
	}

	defer r.Close()

	parse_opts := parser.DefaultParseOptions()
	parse_opts.Body = false

	fm, _, err := parser.ParseReaderWithURI(ctx, parse_opts, r, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s, %s", uri, err)
	}

	return fm, nil
}
