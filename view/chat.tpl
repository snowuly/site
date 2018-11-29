{{define "head"}}
<style>

	.main {
		display: flex;
		flex: 1;
		margin-bottom: 18px;
	}
	.main aside {
		width: 150px;
		background: green;
	}
	.main section {
		flex: 1;
		display: flex;
		flex-direction: column;
		background-color: #red;
	}
	.main section article {
		flex: 1;
		background-color: #888;
		display: flex;
		overflow-y: auto;
		flex-direction: column-reverse;
	}
	.main section footer {
		height: 100px;
		background-color: #777;
		display: flex;
	}
	#output {
		padding: 5px;
	}
	#input {
		height: 100%;
		box-sizing: border-box;
		flex: 1;
	}
	#send {
		width: 100px;		
	}
</style>
{{end}}
{{define "body"}}
{{template "topbar.tpl" .}}
<div class="main">
	<aside>
		<h3>Users List:</h3>
		<div id="list"></div>
	</aside>
	<section>
		<article>
			<div id="output">
			</div>
		</article>
		<footer>
			<textarea disabled id="input"></textarea>
			<button disabled id="send">Send</button>
		</footer>
	</section>
</div>
<script>
	const es = new EventSource("/chat_sse");
	es.onmessage = e => {
		const msgs = e.data.split('\x1f');

		msgs.forEach(msg => {
			switch(true) {
				case msg.startsWith('info:'):
					show(`System: ${msg.slice(5)}`);
					break;
				case msg.startsWith('msg:'):
					show(msg.slice(4));
					break;
				case msg.startsWith('list:'):
					updateList(msg.slice(5).split('\x00'));
					break;
				default:
					console.warn(`unknow msg: ${msg}`);
			}
		});
	}
	es.onerror = e => {
		es.close();
		input.disabled = true
		send.disabled = true
		setTimeout(() => {
			alert("disconnected", e);
		});
	}
	es.onopen = () => {
		input.disabled = false;
		send.disabled = false;
	};

	const output = document.getElementById('output')
	const list = document.getElementById('list')
	const send = document.getElementById('send')
	const input = document.getElementById('input')

	send.addEventListener("click", e => {
		if (input.value === '') {
			return;
		}
		input.disabled = true;
		fetch('/chat_msg', {
			method: 'POST',
			headers: {
				"Content-Type": "application/x-www-form-urlencoded",
			},
			body: `msg=${encodeURIComponent(input.value)}`,
		}).then(res => {
			if (res.status === 200) {
				input.value = '';
				input.disabled = false;
			} else {
				console.log(`send message error: ${res.status}: ${res.statusText}`);
			}
		})
	}, false);

	const el = document.createElement('div');
	el.className = "msg";

	const user = document.createElement('div')
	function updateList(names) {
	console.log(names.length)
		list.innerHTML = '';
		frag = document.createDocumentFragment()
		names.forEach(name => {
			e = user.cloneNode()
			e.textContent = name;
			frag.appendChild(e)
		})
		list.appendChild(frag)
	}
	function show(content) {
		node = el.cloneNode()
		node.textContent = content;
		output.appendChild(node)
	}
</script>
{{end}}
{{template "layout" .}}
