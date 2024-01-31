{{define "rollup" -}}
<ul class="rollup">
{{ range $r := .Rollup -}}
<li><a href="{{ prune_string $r }}/">{{ $r }}</a></li>
{{ end -}}
</ul>
{{ end -}}