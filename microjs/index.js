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

function listenCookieChange(callback, interval = 1000) {
  let lastCookie = document.cookie;
  setInterval(() => {
    let cookie = document.cookie;
    if (cookie !== lastCookie) {
      try {
        callback({ oldValue: lastCookie, newValue: cookie });
      } finally {
        lastCookie = cookie;
      }
    }
  }, interval);
}

function get(path, namespace, params, callback) {
  if (getCookie("micro_refresh")) {
    var expiry = new Date(parseInt(getCookie("micro_expiry")) * 1000);
    if (expiry - Date.now() < 60 * 1000) {
      postCall(
        "auth/Auth/Token",
        "micro",
        {
          refreshToken: getCookie("micro_refresh"),
          options: {
            namespace: namespace,
          },
        },
        function (rsp, status) {
          if (status < 300) {
            setCookie("micro_access", rsp.token.access_token, 30);
            setCookie("micro_refresh", rsp.token.refresh_token, 30);
            setCookie("micro_expiry", rsp.token.expiry, 30);
          }
          getCall(path, namespace, params, callback);
        },
        false
      );
    }
  }
  getCall(path, namespace, params, callback, true)
}

function getCall(path, namespace, params, callback, useToken) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.onreadystatechange = function () {
    if (xmlHttp.readyState == 4);
    callback(JSON.parse(xmlHttp.responseText), xmlHttp.status);
  };
  xmlHttp.open(
    "GET",
    "https://api.m3o.dev/" + path + formatParams(params),
    true
  );
  xmlHttp.setRequestHeader("micro-namespace", namespace);
  if (useToken && getCookie("micro_access")) {
    xmlHttp.setRequestHeader(
      "authorization",
      "Bearer " + getCookie("micro_access")
    );
  }
  xmlHttp.send(null);
}

function post(path, namespace, params, callback) {
  if (getCookie("micro_refresh")) {
    var expiry = new Date(parseInt(getCookie("micro_expiry")) * 1000);
    if (expiry - Date.now() < 60 * 1000) {
      postCall(
        "auth/Auth/Token",
        "micro",
        {
          refreshToken: getCookie("micro_refresh"),
          options: {
            namespace: namespace,
          },
        },
        function (rsp, status) {
          if (status < 300) {
            setCookie("micro_access", rsp.token.access_token, 30);
            setCookie("micro_refresh", rsp.token.refresh_token, 30);
            setCookie("micro_expiry", rsp.token.expiry, 30);
          }
          postCall(path, namespace, params, callback);
        },
        false
      );
    }
  }
  postCall(path, namespace, params, callback, true);
}

function postCall(path, namespace, params, callback, useToken) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.onreadystatechange = function () {
    if (xmlHttp.readyState == 4);
    callback(JSON.parse(xmlHttp.responseText), xmlHttp.status);
  };
  xmlHttp.open("POST", "https://api.m3o.dev/" + path, true); // true for asynchronous
  xmlHttp.setRequestHeader("micro-namespace", namespace);
  if (useToken && getCookie("micro_access")) {
    xmlHttp.setRequestHeader(
      "authorization",
      "Bearer " + getCookie("micro_access")
    );
  }
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
    document.getElementById("errors").style.display = "none";
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
      function (rsp, status) {
        if (status > 299) {
          document.getElementById("errors").innerHTML = rsp.Detail;
          document.getElementById("errors").style.display = "block";
          return;
        }

        if (!rsp && !rsp.token) {
          console.log("Response doesn't look right");
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

  var registerButton = document.getElementById("registerButton");
  registerButton.onclick = function () {
    document.getElementById("errors").style.display = "none";
    Micro.post(
      "signup/sendVerificationEmail",
      "backend",
      {
        email: document.getElementById("emailInput").value,
      },
      function (rsp, status) {
        if (status > 299) {
          document.getElementById("errors").innerHTML = rsp.Detail;
          document.getElementById("errors").style.display = "block";
          return;
        }

        document.getElementById("loginSection").style.display = "none";
        document.getElementById("registerSection").style.display = "none";
        document.getElementById("verifySection").style.display = "block";
        document.getElementById("loginSwitcherSection").style.display = "none";
      }
    );
  };

  var verifyButton = document.getElementById("verifyButton");
  verifyButton.onclick = function () {
    document.getElementById("errors").style.display = "none";
    Micro.post(
      "signup/completeSignup",
      "backend",
      {
        email: document.getElementById("emailInput").value,
        secret: document.getElementById("passwordInput").value,
        token: document.getElementById("verifyInput").value,
        namespace: "backend",
      },
      function (rsp, status) {
        if (status > 299) {
          document.getElementById("errors").innerHTML = rsp.Detail;
          document.getElementById("errors").style.display = "block";
          return;
        }
        if (!rsp && !rsp.token) {
          console.log("Response doesn't look right");
          return;
        }
        setCookie("micro_access", rsp.authToken.access_token, 30);
        setCookie("micro_refresh", rsp.authToken.refresh_token, 30);
        setCookie("micro_expiry", rsp.authToken.expiry, 30);
        // hide modal display if things are all good
        modal.style.display = "none";

        // @todo handle errors
        document.getElementById("loginSection").style.display = "none";
        document.getElementById("registerSection").style.display = "none";
        document.getElementById("verifySection").style.display = "block";
      }
    );
  };

  var signupSwitcher = document.getElementById("signupSwitcher");
  signupSwitcher.onclick = function () {
    document.getElementById("loginSection").style.display = "none";
    document.getElementById("registerSection").style.display = "block";
    document.getElementById("loginSwitcherSection").style.display = "block";
    document.getElementById("signupSwitcherSection").style.display = "none";
    document.getElementById("forgotPassword").style.display = "none";
  };

  var loginSwitcher = document.getElementById("loginSwitcher");
  loginSwitcher.onclick = function () {
    document.getElementById("loginSection").style.display = "block";
    document.getElementById("registerSection").style.display = "none";
    document.getElementById("loginSwitcherSection").style.display = "none";
    document.getElementById("signupSwitcherSection").style.display = "block";
    document.getElementById("forgotPassword").style.display = "block";
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

  Micro["requireLogin"] = function (cb) {
    if (!getCookie("micro_refresh")) {
      btn.click();
      listenCookieChange(({ oldValue, newValue }) => {
        if (newValue) {
          cb();
        }
      }, 200);
    } else {
      cb();
    }
  };
}

initModal();
