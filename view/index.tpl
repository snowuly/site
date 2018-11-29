{{define "head" -}}
<title>Index</title>
{{end}}
{{define "body" -}}
{{template "topbar.tpl" .}}
<ul>
	<li><a href="/chat">Chatroom</a></li>
</ul>
{{end}}
{{- template "layout" .}}
