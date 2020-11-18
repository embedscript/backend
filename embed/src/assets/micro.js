var Micro = (function () {
  "use strict";

  function formatParams(params) {
    return (
      "?" +
      Object.keys(params)
        .map(function (key) {
          return key + "=" + encodeURIComponent(params[key]);
        })
        .join("&")
    );
  }

  // Create the methods object
  var methods = {
    get: function (path, namespace, params, callback) {
      var xmlHttp = new XMLHttpRequest();
      xmlHttp.onreadystatechange = function () {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200);
        callback(JSON.parse(xmlHttp.responseText));
      };
      xmlHttp.open(
        "GET",
        "https://api.m3o.dev/" + path + formatParams(params),
        true
      ); // true for asynchronous
      xmlHttp.setRequestHeader("micro-namespace", namespace);
      xmlHttp.send(null);
    },
  };

  // Expose the public methods
  return methods;
})();
