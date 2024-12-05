GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	rm -rf bin/*
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-mdparse cmd/wof-mdparse/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2feed cmd/wof-md2feed/main.go	
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2html cmd/wof-md2html/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2idx cmd/wof-md2idx/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-md2ts cmd/wof-md2ts/main.go

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
	@make test-prune
	@make test-posts
	@make test-indices
	@make test-feeds

test-prune:
	find fixtures/blog -type f -name '*.html' | xargs rm
	find fixtures/blog -type d -empty -delete

test-posts:
	go run -mod $(GOMOD) cmd/wof-md2html/main.go \
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
	go run -mod $(GOMOD) cmd/wof-md2idx/main.go \
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
	go run -mod $(GOMOD) cmd/wof-md2idx/main.go \
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
	go run -mod $(GOMOD) cmd/wof-md2idx/main.go \
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
	go run -mod $(GOMOD) cmd/wof-md2idx/main.go \
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
	go run -mod $(GOMOD) cmd/wof-md2feed/main.go \
		-format atom_10 \
		-feeds-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-template-uri cwd:///fixtures/templates/feeds \
		blog/

test-feeds-rss:
	go run -mod $(GOMOD) cmd/wof-md2feed/main.go \
		-format rss_20 \
		-feeds-bucket-uri cwd:///fixtures/ \
		-markdown-bucket-uri cwd:///fixtures/ \
		-template-uri cwd:///fixtures/templates/feeds \
		blog/
