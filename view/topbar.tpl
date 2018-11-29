<style>
	header {
		text-align: right;
	}
</style>
<header>
{{ if .IsLogin -}}
Hello {{.Name}}! | <a href="/logout">Logout</a>
{{else -}}
<a href="/login">Login</a> | <a href="/register">Register</a>
{{- end}}
</header>
