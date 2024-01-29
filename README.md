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