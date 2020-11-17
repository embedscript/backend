---
weight: 11
title: Posts
---

# Posts - embed posts in your static site

Transform your static brochure site into an actual blog.

## Get Posts

Display blogs.
JS Fiddle https://jsfiddle.net/rohqmnc1/23/.

```html
<div id="content"></div>
<!-- @todo we should move this to a single micro embeddable js file -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/mustache.js/3.0.0/mustache.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function(event) {
    // @todo move this (or rather, a generic version of it) function to an includable Micro js file
    function getPosts(callback) {
      var xmlHttp = new XMLHttpRequest();
      xmlHttp.onreadystatechange = function() {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
          console.log(xmlHttp.responseText);
        callback(xmlHttp.responseText);
      }
      xmlHttp.open("GET", 'https://api.m3o.dev/posts/query', true); // true for asynchronous 
      xmlHttp.setRequestHeader("micro-namespace", "concert-celtic-uncover")
      xmlHttp.send(null);
    }

    var template = '{{#posts}}<h1>{{title}}</h1>{{/posts}}';

    getPosts(function(data) {
      var result = Mustache.render(template, JSON.parse(data));
      document.getElementById("content").innerHTML = result;
    })
  });

</script>
```
