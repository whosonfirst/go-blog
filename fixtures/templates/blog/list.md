{{define "list" -}}
{{ if ne .Title "" }}
{{ if eq .Mode "tags" }}
## Blog posts tagged <span class="hey-look">{{ .Title }}</span>
{{ else if eq .Mode "authors" }}
## Blog posts written by <span class="hey-look">{{ .Title }}</span>
{{ end }}
{{ end }}
{{ range $fm := .Posts }}
### [{{ $fm.Title }}]({{ $fm.Permalink }})

<a href="{{ $fm.Permalink }}"><img src="{{ $fm.Permalink }}{{ $fm.Image }}" alt="{{ $fm.Title }}" loading="lazy" /></a>

> {{ $fm.Excerpt }}

{{$lena := len $fm.Authors }}
{{$lent := len $fm.Tags }}
<small class="this-is">This is a blog post by
    {{ range $ia, $a := $fm.Authors }}{{ if gt $lena 1 }}{{if eq $ia 0}}{{else if eq (plus1 $ia) $lena}} and {{else}}, {{end}}{{ end }}[{{ $a }}](/blog/authors/{{ prune_string $a  }}){{ end }}.
    {{ if $fm.Date }}It was published on <span class="pubdate"><a href="/blog/{{ $fm.Date.Year }}/{{ $fm.Date.Format "01" }}/">{{ $fm.Date.Format "January" }}</a> <a href="/blog/{{ $fm.Date.Year }}/{{ $fm.Date.Format "01" }}/{{ $fm.Date.Format "02" }}/">{{ $fm.Date.Format "02"}}</a>, <a href="/blog/{{ $fm.Date.Year }}/">{{ $fm.Date.Format "2006" }}</a></span>{{ if gt $lent 0 }} and tagged {{ range $it, $t := $fm.Tags }}{{ if gt $lent 1 }}{{if eq $it 0}}{{else if eq (plus1 $it) $lent}} and {{else}}, {{end}}{{ end }}[{{ $t }}](/blog/tags/{{ prune_string $t  }}){{ end }}{{ end}}.
    {{ else }}
    It was tagged {{ range $it, $t := $fm.Tags }}{{ if gt $lent 1 }}{{if eq $it 0}}{{else if eq (plus1 $it) $lent}} and {{else}}, {{end}}{{ end }}[{{ $t }}](/blog/tags/{{ prune_string $t  }}){{ end }}.
    {{ end }}
</small>
{{ end }}
{{ end }}