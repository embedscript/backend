The `Blogging Embed` turns your static site into a dynamic blog with a few widgets. Lists posts, get a single post and edit them.
<br /><br />

# Get posts

```html
<div id="content"></div>

<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.0/mustache.js"></script>
<script src="https://determined-shaw-741e5d.netlify.app/assets/micro.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function(event) {
    var template = '{{#posts}}<h1><a href="/post?id={{id}}">{{title}}<a/></h1>{{/posts}}';

    Micro.get("posts/query", "concert-celtic-uncover", {}, function(data) {
      var result = Mustache.render(template, data);
      document.getElementById("content").innerHTML = result;
    })
  });
</script>
```
<br /><br / >
# Get single post

Get a single post. The below snippet uses the query parameters of your page to load a post by ID, ie. the `?id=someID` part of the url.

```html
<div id="content"></div>

<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.0/mustache.js"></script>
<script src="https://determined-shaw-741e5d.netlify.app/assets/micro.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function(event) {
    var template = '{{#posts}}<h1>{{title}}</h1>{{/posts}}';

    Micro.get("posts/query", "concert-celtic-uncover", {"id":Micro.params()["id"]}, function(data) {
      var result = Mustache.render(template, data);
      document.getElementById("content").innerHTML = result;
    })
  });
</script>
```