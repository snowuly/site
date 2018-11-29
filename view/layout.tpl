{{define "layout" -}}
<!doctype html>
<html>
	<head>
		<style>
			body {
				display: flex;
				margin: 0;
				flex-direction: column;
				height: 100vh;
			}
			body > .footer {
				text-align: center;
				position: fixed;
				width: 100%;
				bottom: 0;
			}
			
		</style>
		{{template "head" .}}
	</head>
	<body>
		{{template "body" .}}
		{{block "footer" . }}
		<div class="footer">Powerd by <a href="https://github.com/snowuly/kob-go">kob</a></div>
		{{end}}
	</body>
</html>
{{end}}
