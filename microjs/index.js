import * as mustache from './mustache';

var constructor = function () {
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

  function getSearchParameters() {
    var prmstr = window.location.search.substr(1);
    return prmstr != null && prmstr != "" ? transformToAssocArray(prmstr) : {};
  }

  function transformToAssocArray(prmstr) {
    var params = {};
    var prmarr = prmstr.split("&");
    for (var i = 0; i < prmarr.length; i++) {
      var tmparr = prmarr[i].split("=");
      params[tmparr[0]] = tmparr[1];
    }
    return params;
  }

  function getCookieValue(a) {
    var b = document.cookie.match("(^|;)\\s*" + a + "\\s*=\\s*([^;]+)");
    return b ? b.pop() : "";
  }

  var methods = {
    // get makes a get request to the Micro backend
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

    isLoggedIn() {},

    // params returns the query parameters of the current page as an map
    // ie. example.com?a=1&b=2 becomes {"a":"1","b":2"}
    params: getSearchParameters,

    render: mustache.render,
  };

  // Expose the public methods
  return methods;
};

export var Micro = constructor()

