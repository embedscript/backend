var template = "<style>\n  /* The Modal (background) */\n  .modal {\n    display: none; /* Hidden by default */\n    position: fixed; /* Stay in place */\n    z-index: 1; /* Sit on top */\n    padding-top: 100px; /* Location of the box */\n    left: 0;\n    top: 0;\n    width: 100%; /* Full width */\n    height: 100%; /* Full height */\n    overflow: auto; /* Enable scroll if needed */\n    background-color: rgb(0, 0, 0); /* Fallback color */\n    background-color: rgba(0, 0, 0, 0.4); /* Black w/ opacity */\n  }\n\n  /* Modal Content */\n  .modal-content {\n  }\n\n  /* The Close Button */\n  .close {\n    color: #aaaaaa;\n    float: right;\n    font-size: 28px;\n    font-weight: bold;\n  }\n\n  .close:hover,\n  .close:focus {\n    color: #000;\n    text-decoration: none;\n    cursor: pointer;\n  }\n\n  /*\n   Below css copied from https://www.codingnepalweb.com/2020/07/popup-login-form-design-in-html-css.html\n   */\n  @import url(\"https://fonts.googleapis.com/css?family=Poppins:400,500,600,700&display=swap\");\n\n  .show-btn {\n    background: #fff;\n    padding: 10px 20px;\n    font-size: 20px;\n    font-weight: 500;\n    color: #3498db;\n    cursor: pointer;\n    box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);\n  }\n  .show-btn,\n  .container {\n    position: absolute;\n    top: 50%;\n    left: 50%;\n    transform: translate(-50%, -50%);\n  }\n\n  input[type=\"checkbox\"] {\n    display: none;\n  }\n  .container {\n    background: #fff;\n    width: 410px;\n    padding: 30px;\n    box-shadow: 0 0 8px rgba(0, 0, 0, 0.1);\n  }\n  #show:checked ~ .container {\n    display: block;\n  }\n  .container .close-btn {\n    position: absolute;\n    right: 20px;\n    top: 15px;\n    font-size: 18px;\n    cursor: pointer;\n  }\n  .container .close-btn:hover {\n    color: #3498db;\n  }\n  .container .text {\n    font-size: 35px;\n    font-weight: 600;\n    text-align: center;\n  }\n  .container form {\n    margin-top: -20px;\n  }\n  .container form .data {\n    height: 45px;\n    width: 100%;\n    margin: 40px 0;\n  }\n  form .data label {\n    font-size: 18px;\n  }\n  form .data input {\n    height: 100%;\n    width: 100%;\n    padding-left: 10px;\n    font-size: 17px;\n    border: 1px solid silver;\n  }\n  form .data input:focus {\n    border-color: #3498db;\n    border-bottom-width: 2px;\n  }\n  form .forgot-pass {\n    margin-top: -8px;\n  }\n  form .forgot-pass a {\n    color: #3498db;\n    text-decoration: none;\n  }\n  form .forgot-pass a:hover {\n    text-decoration: underline;\n  }\n  form .btn {\n    margin: 30px 0;\n    height: 45px;\n    width: 100%;\n    position: relative;\n    overflow: hidden;\n  }\n  form .btn .inner {\n    height: 100%;\n    width: 300%;\n    position: absolute;\n    left: -100%;\n    z-index: -1;\n    background: -webkit-linear-gradient(\n      right,\n      #56d8e4,\n      #9f01ea,\n      #56d8e4,\n      #9f01ea\n    );\n    transition: all 0.4s;\n  }\n  form .btn:hover .inner {\n    left: 0;\n  }\n  form .btn button {\n    height: 100%;\n    width: 100%;\n    background: none;\n    border: none;\n    color: #fff;\n    font-size: 18px;\n    font-weight: 500;\n    text-transform: uppercase;\n    letter-spacing: 1px;\n    cursor: pointer;\n  }\n  form .signup-link {\n    text-align: center;\n  }\n  form .signup-link a {\n    color: #3498db;\n    text-decoration: none;\n  }\n  form .signup-link a:hover {\n    text-decoration: underline;\n  }\n</style>\n\n<!-- Trigger/Open The Modal -->\n<button id=\"myBtn\" style=\"display: none\">Open Modal</button>\n\n<!-- The Modal -->\n<div id=\"myModal\" class=\"modal\">\n  <!-- Modal content -->\n  <div class=\"modal-content\">\n    <span style=\"display: none\" class=\"close\">&times;</span>\n    <!-- Below popup content was copied from https://www.codingnepalweb.com/2020/07/popup-login-form-design-in-html-css.html -->\n    <div class=\"center\">\n      <input type=\"checkbox\" id=\"show\">\n      <label for=\"show\" class=\"show-btn\">View Form</label>\n      <div class=\"container\">\n        <label for=\"show\" class=\"close-btn fas fa-times\" title=\"close\"></label>\n        <div class=\"text\">Login Form</div>\n        <form action=\"#\">\n          <div class=\"data\">\n            <label>Email or Phone</label>\n            <input id=\"emailInput\" type=\"text\" required>\n          </div>\n          <div class=\"data\">\n            <label>Password</label>\n            <input id=\"passwordInput\" type=\"password\" required>\n          </div>\n          <div class=\"forgot-pass\">\n            <a href=\"#\">Forgot Password?</a>\n          </div>\n          <div class=\"btn\">\n            <div class=\"inner\"></div>\n            <button id=\"loginButton\" type=\"button\">login</button>\n          </div>\n          <div class=\"signup-link\">\n            Not a member? <a href=\"#\">Signup now</a>\n          </div>\n        </form>\n      </div>\n    </div>\n  </div>\n</div>\n";

//import * as mustache from "./mustache";

function formatParams(params) {
  return "?" + Object.keys(params).map(function (key) {
    return key + "=" + encodeURIComponent(params[key]);
  }).join("&");
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

function get(path, namespace, params, callback) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.onreadystatechange = function () {
    if (xmlHttp.readyState == 4 && xmlHttp.status == 200) ;
    callback(JSON.parse(xmlHttp.responseText));
  };
  xmlHttp.open("GET", "https://api.m3o.dev/" + path + formatParams(params), true); // true for asynchronous
  xmlHttp.setRequestHeader("micro-namespace", namespace);
  xmlHttp.send(null);
}

function post(path, namespace, params, callback) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.onreadystatechange = function () {
    if (xmlHttp.readyState == 4 && xmlHttp.status == 200) ;
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
    params: getSearchParameters
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
    Micro.post("auth/Auth/Token", "micro", {
      id: document.getElementById("emailInput").value,
      secret: document.getElementById("passwordInput").value,
      options: {
        namespace: "backend"
      }
    }, function (text) {
      console.log(text);
    });
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