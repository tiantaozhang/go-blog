<h1>Blog posts</h1>
<p>
{{.title}}	
</p>
<div>
<ul>

    {{range .blogs}}
    {{with .}}
    <li>
		<a href="/view/{{._id}}">{{.title}}</a>
        <a href="/edit/{{._id}}">Edit</a>
        <a href="/delete/{{._id}}">Delete</a>
    </li>
    {{end}}
    {{end}}
</ul>

</div>

     

