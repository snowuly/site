{{define "head" -}}
<title>Register</title>
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
		<h3>Register</h3>
		<form method="POST" action="/register">
			<input name="login" placeholder="login name" />
			<input name="nickname" placeholder="nickname" />
			<input name="pwd"type="password" placeholder="password"/>
			<input name="repwd" type="password" placeholder="confirm password"/>
			<input type="submit">
		</form>
	</div>
</div>
{{end}}
{{template "layout" .}}
