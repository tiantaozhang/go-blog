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
