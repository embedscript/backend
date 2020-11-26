The `Blogging Embed` turns your static site into a dynamic blog with a few widgets. Lists posts, get a single post and edit them.
<br /><br />

# Get posts

```html
<div id="content"></div>

<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.0/mustache.js"></script>
<script src="https://determined-shaw-741e5d.netlify.app/assets/micro.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function (event) {
    var template =
      '{{#posts}}<h2><a href="/post?id={{id}}">{{title}}<a/></h2>{{/posts}}';

    Micro.post(
      "posts/query",
      "backend",
      {
        website: "example.com",
      },
      function (data) {
        var result = Mustache.render(template, data);
        document.getElementById("content").innerHTML = result;
      }
    );
  });
</script>
```
[JSFiddle]()

<br /><br />
# Get single post

Get a single post. The below snippet uses the query parameters of your page to load a post by ID, ie. the `?id=someID` part of the url.

```html
<div id="content"></div>

<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.0/mustache.js"></script>
<script src="https://determined-shaw-741e5d.netlify.app/assets/micro.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function (event) {
    var template =
      '{{#posts}}<h2><a href="/post?id={{id}}">{{title}}<a/></h2>{{/posts}}';

    Micro.post(
      "posts/query",
      "backend",
      {
        website: "example.com",
        id: Micro.params()["postID"],
      },
      function (data) {
        var result = Mustache.render(template, data);
        document.getElementById("content").innerHTML = result;
      }
    );
  });
</script>
```

# Save post

```html
<div id="content"></div>

<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.0/mustache.js"></script>
<script src="https://determined-shaw-741e5d.netlify.app/assets/micro.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function (event) {
    var template = `
      <input id="post-title"><br /><br />
      <textarea id="post-content"></textarea><br /><br />
      <button id="save-button" type="button">Save</button>
    `;
    document.getElementById("content").innerHTML = template;

    document.getElementById("save-button").onclick = function () {
      Micro.requireLogin(function () {
        Micro.post(
          "posts/save",
          "backend",
          {
            website: "example.com",
            title: document.getElementById("post-title").value,
            content: document.getElementById("post-content").value,
          },
          function (data) {
            console.log("Successfully saved.");
          }
        );
      });
    };
  });
</script>
```
