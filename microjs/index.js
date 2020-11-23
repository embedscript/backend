"use strict";
//import * as mustache from "./mustache";
import template from "./popup.html";

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

function setCookie(name, value, days) {
  var expires = "";
  if (days) {
    var date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(";");
  for (var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == " ") c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}

function eraseCookie(name) {
  document.cookie = name + "=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;";
}

function get(path, namespace, params, callback) {
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
}

function post(path, namespace, params, callback) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.onreadystatechange = function () {
    if (xmlHttp.readyState == 4 && xmlHttp.status == 200);
    callback(JSON.parse(xmlHttp.responseText));
  };
  xmlHttp.open("POST", "https://api.m3o.dev/" + path, true); // true for asynchronous
  xmlHttp.setRequestHeader("micro-namespace", namespace);
  xmlHttp.send(JSON.stringify(params));
}

var constructor = function () {
  var methods = {
    // get makes a get request to the Micro backend
    get: get,
    post: post,
    isLoggedIn() {
      var v = getCookieValue("micro_session");
      if (v.length == 0) {
        return false;
      }
      return true;
    },
    // params returns the query parameters of the current page as an map
    // ie. example.com?a=1&b=2 becomes {"a":"1","b":2"}
    params: getSearchParameters,
    // render: mustache.render,
  };

  // Expose the public methods
  return methods;
};

var Micro = constructor();

function initModal() {
  document.body.innerHTML += template;

  // Set up handlers for modal
  // Get the modal
  var modal = document.getElementById("myModal");

  // Get the button that opens the modal
  var btn = document.getElementById("myBtn");

  // Get the <span> element that closes the modal
  var span = document.getElementsByClassName("close")[0];

  // Get the button that opens the modal
  var loginButton = document.getElementById("loginButton");
  loginButton.onclick = function () {
    Micro.post(
      "auth/Auth/Token",
      "micro",
      {
        id: document.getElementById("emailInput").value,
        secret: document.getElementById("passwordInput").value,
        options: {
          namespace: "backend",
        },
      },
      function (rsp) {
        if (!rsp && !rsp.token) {
          return;
        }
        setCookie("micro_access", rsp.token.access_token, 30);
        setCookie("micro_refresh", rsp.token.refresh_token, 30);
        setCookie("micro_expiry", rsp.token.expiry, 30);
        // hide modal display if things are all good
        modal.style.display = "none";
      }
    );
  };

  // When the user clicks the button, open the modal
  btn.onclick = function () {
    modal.style.display = "block";
  };

  // When the user clicks on <span> (x), close the modal
  span.onclick = function () {
    modal.style.display = "none";
  };

  // When the user clicks anywhere outside of the modal, close it
  window.onclick = function (event) {
    if (event.target == modal) {
      modal.style.display = "none";
    }
  };

  Micro["popup"] = function () {
    btn.click();
  };
}

initModal();
