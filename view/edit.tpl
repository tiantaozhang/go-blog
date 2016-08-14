<h1>Edit {{.title}}</h1>

<form action="" method="post">
    标题:<input id="t1" type="text" name="title" value="{{.title}}"><br>
    内容：<textarea name="content" colspan="30" rowspan="50">{{.content}}</textarea>
    <input type="hidden" name="id" value="{{._id}}">
    <input type="submit">
</form>