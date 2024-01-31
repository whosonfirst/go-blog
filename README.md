# go-blog

There are many "static-site" blogging tools. This one is ours.

## Motivation

This is designed to be a very simple blogging system and static-site generators. Posts are written in Markdown using Jekyll-style "front matter" blocks with support for a handful of custom properties.

Those posts are then published as HTML in to /YYYY/MM/DD/POSTNAME directory trees. Index-style pages, with slugs, can be generated for all the posts in reverse-chronological order as well as the authors, tags and dates associated with posts. Date indices are generated for year, year-month and year-month-day combinations. Additionally, Atom 1.0 and/or RSS 2.0 syndication feeds can be generated for the most recent posts.

That's it. In many ways this is just a simplified version of Jekyll, written in Go, but that's sort of the point. While it may a little bit of time to set up custom templates once that's done these tools are designed to "just work" more or less forever as-is. As such there's a bunch of stuff these tools don't do, or don't do yet, like hidden drafts or workflow processes.

For a concrete example have a look at the [fixtures/blog](fixtures/blog) folder and the `test-*` targets in the [Makefile](Makefile).

## FrontMatter

Posts should start with a Jekyll-style "frontmatter" block with the following keys (and values updated to reflect the post being written):

```
---
layout: page
permalink: /blog/2024/01/28/test/
published: true
title: This is a test
date: 2024-01-28
category: blog
excerpt: "This test left intentionally blank"
authors: [author1,author2]
image: images/test.jpg
tags: [tag1,tag2,tag3]
---
```

## Templates

### wof-md2html

The `wof-md2html` does not require any user-defined templates but it does support custom "header" and "footer" templates to wrap the output of any given post.

Examples:

* [fixtures/templates/blog/header.html](fixtures/templates/blog/header.html)
* [fixtures/templates/blog/footer.html](fixtures/templates/blog/footer.html)

### wof-md2idx

The `wof-md2idx` does not require any user-defined templates but it does support custom "header" and "footer" templates to wrap the output of any given index, as well a "list" template for list views and a "rollup" template for rollup views (for example, all the tags or all the authors)..

Examples:

* [fixtures/templates/blog/header.html](fixtures/templates/blog/header.html)
* [fixtures/templates/blog/footer.html](fixtures/templates/blog/footer.html)
* [fixtures/templates/blog/list.html](fixtures/templates/blog/list.html)
* [fixtures/templates/blog/rollup.html](fixtures/templates/blog/rollup.html)

### wof-md2feed

If you are going to use the `wof-md2feed` tool you will need to ensure that templates with the following names are loaded:

#### feed_atom_10

A Go language template defining an Atom 1.0 syndication feed.

Example: [fixtures/templates/feeds/atom_10.xml](fixtures/templates/feeds/atom_10.xml)

#### feeds_rss_20

A Go language template defining an RSS 2.0 syndication feed.

Example: [fixtures/templates/feeds/rss_20.xml](fixtures/templates/feeds/rss_20.xml)

## "Buckets"

Under the hood the code uses the [GoCloud `Blob` abstraction layer](https://gocloud.dev/howto/blob/) for reading and writing files. These include source Markdown files, HTML files that are generated and any (Go language) template files used to supplement or decorate the default HTML output.

As of this writing only the [file://](https://gocloud.dev/howto/blob/#local) protocol handler is supported by default in addition to the non-standard helper protocol `cwd://` which will attempt to derive a `file://` URI for the current working directory.

For example, if you were in a directory called `/usr/local/weblog` then the URI `cwd://` would be interpreted as `file:///usr/local/weblog`.

In future releases other `Blob` providers, notably S3, will be supported.

## Tools

```
$> make cli
rm -rf bin/*
go build -mod vendor -ldflags="-s -w" -o bin/wof-mdparse cmd/wof-mdparse/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-md2feed cmd/wof-md2feed/main.go	
go build -mod vendor -ldflags="-s -w" -o bin/wof-md2html cmd/wof-md2html/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-md2idx cmd/wof-md2idx/main.go
```

### wof-mdparse

Parse one or more whosonfirst/go-blog -style Markdown URIs and output FrontMatter, body text or both.

```
$> ./bin/wof-mdparse -h
Parse one or more whosonfirst/go-blog -style Markdown URIs and output FrontMatter, body text or both.
Usage of ./bin/wof-mdparse:
  -all
    	Dump both frontmatter and body
  -body
    	Dump (Markdown) body
  -frontmatter
    	Dump (Jekyll) frontmatter
  -markdown-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where Markdown files should be read from.
```

For example:

```
$> ./bin/wof-mdparse \
	-frontmatter \
	-markdown-bucket-uri file:///usr/local/sfomuseum/www-sfomuseum-weblog/www/ \
	blog/2024/01/22/shoebox/index.md
	
---
layout: page
permalink: /blog/2024/01/22/shoebox/
published: true
title: The SFO Museum Aviation Collection Website Shoebox
date: 2024-01-22 00:00:00 +0000 UTC
category: blog
excerpt: This is a blog post about something that’s been hiding in plain sight on the SFO Museum Aviation Collection website for over a month now: The ability to save collection objects to a personal “shoebox”. If that sounds like a simple bookmarking system limited to items in the SFO Museum collection that’s because it is. For now. The shoebox and the introduction of user accounts are the first steps, the first building blocks, towards developing more sophisticated functionality and applications for the museum and its collection.
authors: [aaron cope]
image: images/1762892311_inyt2xFHYjWknvCGH6bITzH6doudbSzs_c.jpg
tags: [collection shoebox]
---
```

### wof-md2html

Converts a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and converts them to HTML documents writing them to a target gocloud.dev/blob bucket URI.

```
$> ./bin/wof-md2html -h
Converts a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and converts them to HTML documents writing them to a target gocloud.dev/blob bucket URI.
Usage:
	 ./bin/wof-md2html [options] uri(N) uri(N)
  -footer string
    	The name of the (Go) template to use as a custom footer
  -header string
    	The name of the (Go) template to use as a custom header
  -html-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where HTML files should be written to.
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -markdown-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where Markdown files should be read from.
  -mode string
    	Valid modes are: files, directory (default "files")
  -output string
    	What you expect the output HTML file to be called (default "index.html")
  -template-uri value
    	One or more valid gocloud.dev/blob bucket URIs where HTML template files should be read from.
```

For example:

```
$> ./bin/wof-md2html \
	-markdown-bucket-uri file:///usr/local/sfomuseum/www-sfomuseum-weblog/www/ \
	-html-bucket-uri file:///usr/local/sfomuseum/www-sfomuseum-weblog/www/ \
	blog/2024/01/22/shoebox/index.md
```

### wof-md2idx

Generate paginated "index"-style list pages for a collection of blog posts. List styles include authors, tags, dates and reverse-chronological posts.

```
$> ./bin/wof-md2idx -h
Generate paginated "index"-style list pages for a collection of blog posts. List styles include authors, tags, dates and reverse-chronological posts.
Usage of ./bin/wof-md2idx:
  -footer string
    	The name of the (Go) template to use as a custom footer
  -header string
    	The name of the (Go) template to use as a custom header
  -html-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where HTML files should be written to.
  -html-template-uri value
    	One or more valid gocloud.dev/blob bucket URIs where HTML template files should be read from.
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -list string
    	The name of the (Go) template to use as a custom list view
  -markdown-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where Markdown files should be read from.
  -markdown-template-uri value
    	One or more valid gocloud.dev/blob bucket URIs where Markdown template files should be read from.
  -mode string
    	Valid modes are: authors, landing, tags, ymd. (default "landing")
  -output string
    	What you expect the output HTML file to be called (default "index.html")
  -per-page int
    	The number of posts to include on a single page (the rest will be paginated on to page2, page3 and so on. (default 10)
  -rollup string
    	The name of the (Go) template to use as a custom rollup view (for things like tags and authors)
```

### wof-md2feed

Generate Atom 1.0 or RSS 2.0 syndication feeds from a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and writing the feeds to a target gocloud.dev/blob bucket URI.

```
$> ./bin/wof-md2feed -h
Generate Atom 1.0 or RSS 2.0 syndication feeds from a collection of Markdown documents read from a source gocloud.dev/blob bucket URI and writing the feeds to a target gocloud.dev/blob bucket URI.
Usage:
	 ./bin/wof-md2feed [options] uri(N) uri(N)
  -feeds-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where feeds should be written to.
  -format string
    	Valid options are: atom_10, rss_20 (default "rss_20")
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -items int
    	The number of items to include in your feed (default 10)
  -markdown-bucket-uri string
    	A valid gocloud.dev/blob bucket URI where Markdown files should be read from.
  -output string
    	The filename of your feed. If empty default to the value of -format + ".xml"
  -template-uri value
    	One or more valid gocloud.dev/blob bucket URIs where feed template files should be read from.
```

## Putting it all together

Here are some _example_ Makefile targets for a weblog where copies the binary tools produced by this package are stored in a folder called `dist`, templates for the blog are stored in `templates` and the Markdown files and resultant HTML files are stored in `www/blog`. Note that these targets do not do anything to publish or sync the rendered Markdown files to a remote server or location. Those details are left to you to figure out for yourself.

_For a working example consult the `test*` targets in the [Makefile](Makefile)._

```
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
BIN=utils/$(OS)

# https://github.com/whosonfirst/go-blog
MD2HTML=$(BIN)/wof-md2html
MD2IDX=$(BIN)/wof-md2idx
MD2FEED=$(BIN)/wof-md2feed

render-blog:
	@make render-blog-posts
	@make render-blog-indices
	@make render-blog-feeds

render-blog-posts:
	$(MD2HTML) \
		-template-uri cwd:///templates/common \
		-template-uri cwd:///templates/blog/post \
		-header blog_post_header \
		-footer blog_post_footer \
		-mode directory \
		-html-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/

render-blog-indices:
	@make render-blog-landing
	@make render-blog-ymd
	@make render-blog-tags
	@make render-blog-authors

render-blog-landing:
	$(MD2IDX) -mode landing \
		-html-template-uri cwd:///templates/common \
		-html-template-uri cwd:///templates/blog/index \
		-markdown-template-uri cwd:///templates/blog/index \
		-header blog_index_header \
		-footer blog_index_footer \
		-list blog_index_list \
		-html-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/

render-blog-ymd:
	$(MD2IDX) -mode ymd \
		-html-template-uri cwd:///templates/common \
		-html-template-uri cwd:///templates/blog/index \
		-markdown-template-uri cwd:///templates/blog/index \
		-header blog_index_header \
		-footer blog_index_footer \
		-list blog_index_list \
		-html-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/

render-blog-tags:
	$(MD2IDX) -mode tags \
		-html-template-uri cwd:///templates/common \
		-html-template-uri cwd:///templates/blog/index \
		-markdown-template-uri cwd:///templates/blog/index \
		-header blog_index_header \
		-footer blog_index_footer \
		-list blog_index_list \
		-rollup blog_index_rollup \
		-html-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/

render-blog-authors:
	$(MD2IDX) -mode authors \
		-html-template-uri cwd:///templates/common \
		-html-template-uri cwd:///templates/blog/index \
		-markdown-template-uri cwd:///templates/blog/index \
		-header blog_index_header \
		-footer blog_index_footer \
		-list blog_index_list \
		-rollup blog_index_rollup \
		-html-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/

render-blog-feeds:
	@make render-blog-feeds-rss
	@make render-blog-feeds-atom

render-blog-feeds-rss:
	$(MD2FEED) \
		-template-uri cwd:///templates/blog/feed \
		-format rss_20 \
		-feeds-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/

render-blog-feeds-atom:
	$(MD2FEED) \
		-template-uri cwd:///templates/blog/feed \
		-format atom_10 \
		-feeds-bucket-uri cwd:///www/ \
		-markdown-bucket-uri cwd:///www/ \
		blog/
```

## See also

* https://github.com/russross/blackfriday/v2
* https://gocloud.dev/howto/blob/