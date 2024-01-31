GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	rm -rf bin/*
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-mdparse cmd/wof-mdparse/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2feed cmd/wof-md2feed/main.go	
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2html cmd/wof-md2html/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2idx cmd/wof-md2idx/main.go

dist-build:
	# OS=darwin make dist-os
	OS=windows make dist-os
	OS=linux make dist-os

dist-os:
	mkdir -p dist/$(OS)
	GOOS=$(OS) GOARCH=386 go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o dist/$(OS)/wof-mdparse cmd/wof-mdparse/main.go
	GOOS=$(OS) GOARCH=386 go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o dist/$(OS)/wof-md2feed cmd/wof-md2feed/main.go
	GOOS=$(OS) GOARCH=386 go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o dist/$(OS)/wof-md2html cmd/wof-md2html/main.go
	GOOS=$(OS) GOARCH=386 go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o dist/$(OS)/wof-md2idx cmd/wof-md2idx/main.go

test:
	@make test-posts
	@make test-indices
	@make test-feeds

test-posts:
	bin/wof-md2html \
		-mode directory \
		-html-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-template-uri cwd:///fixtures/templates/blog \
		-header header \
		-footer footer \
		blog/

test-indices:
	@make test-indices-tags
	@make test-indices-authors
	@make test-indices-ymd
	@make test-indices-landing

test-indices-tags:
	bin/wof-md2idx \
		-mode tags \
		-html-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-html-template-uri cwd:///fixtures/templates/blog \
		-header header \
		-footer footer \
		-list list \
		-rollup rollup \
		blog/

test-indices-authors:
	bin/wof-md2idx \
		-mode authors \
		-html-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-html-template-uri cwd:///fixtures/templates/blog \
		-header header \
		-footer footer \
		-list list \
		-rollup rollup \
		blog/

test-indices-ymd:
	bin/wof-md2idx \
		-mode ymd \
		-html-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-html-template-uri cwd:///fixtures/templates/blog \
		-header header \
		-footer footer \
		-list list \
		-rollup rollup \
		blog/

test-indices-landing:
	bin/wof-md2idx \
		-mode landing \
		-html-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-html-template-uri cwd:///fixtures/templates/blog \
		-header header \
		-footer footer \
		-list list \
		-rollup rollup \
		blog/

test-feeds:
	@make test-feeds-atom
	@make test-feeds-rss

test-feeds-atom:
	bin/wof-md2feed \
		-format atom_10 \
		-feeds-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-template-uri cwd:///fixtures/templates/feeds \
		blog/

test-feeds-rss:
	bin/wof-md2feed \
		-format rss_20 \
		-feeds-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-template-uri cwd:///fixtures/templates/feeds \
		blog/
