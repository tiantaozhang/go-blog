
<!DOCTYPE html>
<html>
<head>
	<title>blog</title>
	<meta charset="utf-8">
	<link rel="stylesheet" type="text/css" href="/static/view.css">
</head>
<body>
<div id="main">
	<header id="head">
		<nav id="nav">
		<div class="logo"><a href="#">张家界</a>
		</div>
			<ul>	
				<li><a href="#">登录</a></li>
				<li><a href="#">注册</a></li>
				<li><a href="#">联系我们</a></li>
			</ul>
		</nav>
		<div id="banner">
			<div id="inner">
				<h1>张的主页</h1>
				<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
				tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,
				quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
				</p>
				<button>了解我</button>
			</div>
		</div>
	</header>

	<content id="container">
		<div class="subject">
			
		</div>
		<div class="sidebar">
			<div id="avatar">
				<img src="/static/me.jpg" title="头像">
				<div class="nickname"><a href="#" class="username">tatumn zhang</a></div>
			</div>
			<hr style="border:1px dashed #cec6c6; height: 1px;border-top: none;">
			<div id="visit">
				<ul>
					<li>
						访问量:
						<span></span>
					</li>
					<li>
						今日访问量:
						<span></span>
					</li>
				</ul>			
			</div>
			<hr style="border:1px dashed #cec6c6; height: 1px;border-top: none;">
			<div id="comment">
				
			</div>
			<hr style="border:1px dashed #cec6c6; height: 1px;border-top: none;">
			<div class="readcount">
				
			</div>
		</div>
	
	</content>
	<footer>
		<ul>
			<li>a</li>
			<li>b</li>
			<li>c</li>
			<li>d</li>
		</ul>
		<div id="foot">
			<span>Copyright © Tiantao Zhang</span>
		</div>
	</footer>
</div>
</body>
</html>


<!-- 
<!DOCTYPE html>
<html>
<head>
	<title>view</title>
	<meta charset="utf-8">
	<link rel="stylesheet" type="text/css" href="/static/view.css">
</head>

<body>

<div id="head">
<h1>{{.title}}</h1>	

<div>
	<div class="author">作者:{{.author}}</div>
	<div class="createtime">创建时间:{{.time}}</div>
</div>

</div>

<div id="container">

<div class="main">
	{{.content}}	
</div>

<div class="sidebar">
	
</div>

</div>


<div id="bottom">footer</div>


</body>
</html>
 -->