# go-blog

There are many blogging tools. This one is ours.

## Documentation

Documentation is incomplete at this point.

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

For example:

```
$> ./bin/wof-md2html \
	-markdown-bucket-uri file:///usr/local/sfomuseum/www-sfomuseum-weblog/www/ \
	-html-bucket-uri file:///usr/local/sfomuseum/www-sfomuseum-weblog/www/ \
	blog/2024/01/22/shoebox/index.md
```

## Putting it all together

Here are some _example_ Makefile targets for a weblog where copies the binary tools produced by this package are stored in a folder called `dist`, templates for the blog are stored in `templates` and the Markdown files and resultant HTML files are stored in `www/blog`. Note that these targets do not do anything to publish or sync the rendered Markdown files to a remote server or location. Those details are left to you to figure out for yourself.

```
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
BIN=utils/$(OS)

# https://github.com/whosonfirst/go-whosonfirst-markdown
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