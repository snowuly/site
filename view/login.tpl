
{{define "head" -}}
<title>Login</title>
<style>
	.content {
		display: flex;
		justify-content: center;
		align-items: center;
		box-sizing: border-box;
		height: 80vh;
	}
	.content h3 {
		width: 100%;
	}
	
</style>
{{end}}
{{define "body" -}}
<div class="content">
	<div>
		<h3>Login</h3>
		<form method="POST" action="/login">
			<input name="login" placeholder="login name" />
			<input name="pwd"type="password" placeholder="password"/>
			<input type="submit">
		</form>
		<p>No account? Go to <a href="/register">Register</a></p>
	</div>
</div>
{{end}}
{{template "layout" .}}
