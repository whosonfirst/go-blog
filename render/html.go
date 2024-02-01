package render

import (
	"bytes"
	"html/template"
	"io"
	"log"

	"github.com/russross/blackfriday/v2"
	"github.com/whosonfirst/go-blog/jekyll"
	"github.com/whosonfirst/go-blog/markdown"
	"gocloud.dev/blob"
)

type Pagination struct {
	Page     int
	Pages    int
	Total    int
	Next     string
	Previous string
}

type HTMLOptions struct {
	Mode         string
	Input        string
	Output       string
	Header       string
	Footer       string
	List         string
	Title        string
	Pagination   *Pagination
	Templates    *template.Template
	TargetBucket *blob.Bucket
	SourceBucket *blob.Bucket
	PerPage      int
}

func DefaultHTMLOptions() *HTMLOptions {

	opts := HTMLOptions{
		Mode:      "files",
		Input:     "index.md",
		Output:    "index.html",
		Header:    "",
		Footer:    "",
		List:      "",
		Templates: nil,
	}

	return &opts
}

type nopCloser struct {
	io.Reader
}

type WOFRenderer struct {
	bf          *blackfriday.HTMLRenderer
	frontmatter *jekyll.FrontMatter
	header      string
	footer      string
	templates   *template.Template
	pagination  *Pagination
}

type WOFRendererHeaderVars struct {
	FrontMatter *jekyll.FrontMatter
	Pagination  *Pagination
}

type WOFRendererFooterVars struct {
	FrontMatter *jekyll.FrontMatter
	Pagination  *Pagination
}

func (r *WOFRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	switch node.Type {

	case blackfriday.Image:
		return r.bf.RenderNode(w, node, entering)
	default:
		return r.bf.RenderNode(w, node, entering)
	}
}

func (r *WOFRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {

	if r.templates == nil || r.header == "" {
		r.bf.RenderHeader(w, ast)
		return
	}

	t := r.templates.Lookup(r.header)

	if t == nil {
		log.Printf("Invalid or missing template '%s'\n", r.header)
		return
	}

	vars := WOFRendererHeaderVars{
		FrontMatter: r.frontmatter,
		Pagination:  r.pagination,
	}

	err := t.Execute(w, vars)

	if err != nil {
		log.Println(err)
	}
}

func (r *WOFRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {

	if r.templates == nil || r.footer == "" {
		r.bf.RenderFooter(w, ast)
		return
	}

	t := r.templates.Lookup(r.footer)

	if t == nil {
		log.Printf("Invalid or missing template '%s'\n", r.footer)
		return
	}

	vars := WOFRendererFooterVars{
		FrontMatter: r.frontmatter,
		Pagination:  r.pagination,
	}

	err := t.Execute(w, vars)

	if err != nil {
		log.Println(err)
	}
}

func (nopCloser) Close() error { return nil }

func RenderHTML(d *markdown.Document, opts *HTMLOptions) (io.ReadCloser, error) {

	flags := blackfriday.CommonHTMLFlags
	flags |= blackfriday.CompletePage
	flags |= blackfriday.UseXHTML

	params := blackfriday.HTMLRendererParameters{
		Flags: flags,
	}

	renderer := blackfriday.NewHTMLRenderer(params)

	r := WOFRenderer{
		bf:          renderer,
		frontmatter: d.FrontMatter,
		header:      opts.Header,
		footer:      opts.Footer,
		templates:   opts.Templates,
		pagination:  opts.Pagination,
	}

	unsafe := blackfriday.Run(d.Body.Bytes(), blackfriday.WithRenderer(&r))

	// safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	html := bytes.NewReader(unsafe)
	return nopCloser{html}, nil

}
