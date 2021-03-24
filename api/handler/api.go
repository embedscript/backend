package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	pb "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	filesproto "github.com/micro/services/files/proto"
)

type V1 struct{}

func (e *V1) Serve(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	files := filesproto.NewFilesService("files", client.DefaultClient)

	if req.Get["project"] == nil || len(req.Get["project"].Values) == 0 {
		return errors.New("bad request")
	}
	local := false
	if req.Get["local"] != nil {
		local = true
	}
	project := req.Get["project"].Values[0]
	script := ""
	if req.Get["script"] != nil && len(req.Get["script"].Values) > 0 {
		script = req.Get["script"].Values[0]
		logger.Infof("Serving %v", script)
	}

	htmlFile := ""
	jsFile := ""
	cssFile := ""
	owner := ""

	if req.Get["javascript"] != nil && len(req.Get["javascript"].Values) > 0 &&
		req.Get["html"] != nil && len(req.Get["html"].Values) > 0 &&
		req.Get["css"] != nil && len(req.Get["css"].Values) > 0 {
		jsFile = req.Get["javascript"].Values[0]
		htmlFile = req.Get["html"].Values[0]
		cssFile = req.Get["html"].Values[0]
	} else {
		resp, err := files.List(ctx, &filesproto.ListRequest{
			Project: script,
		})
		if err != nil {
			return err
		}
		if len(resp.Files) == 0 {
			return errors.New("not found")
		}

		for _, file := range resp.Files {
			if file.Owner != "" {
				owner = file.Owner
			}
			switch {
			case strings.Contains(file.Path, "main"):
				jsFile = file.FileContents
			case strings.Contains(file.Path, "index"):
				htmlFile = file.FileContents
			case strings.Contains(file.Path, "style"):
				cssFile = file.FileContents
			}
		}
	}

	id, _ := uuid.NewV4()

	scriptTags, linkTags, html := extractScriptLink(htmlFile)
	htmlFile = html

	prodLinks := map[string]string{
		"microjs":       "https://embedscript.com/assets/micro.js",
		"diff-renderer": "https://cdn.jsdelivr.net/npm/morphdom@2.6.1/dist/morphdom.min.js",
	}
	localLinks := map[string]string{
		"microjs":       "http://127.0.0.1:4200/assets/micro.js",
		"diff-renderer": "https://cdn.jsdelivr.net/npm/morphdom@2.6.1/dist/morphdom.min.js",
	}
	links := prodLinks
	if local {
		links = localLinks
	}

	rendered := `<html>
	<head>
	    ` + linkTags + `
		<style>` +
		cssFile +
		`</style>
	</head>
	<body>
	<div id="` + id.String() + `">
	</div>
	` + scriptTags + `
	<script src="` + links["microjs"] + `"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/handlebars.js/4.7.7/handlebars.min.js"></script>
	<script id="template" type="x-tmpl-mustache">` +
		htmlFile + `
	</script>
	<script type=module>
	import {default as morphdom} from 'https://unpkg.com/morphdom?module';

	Handlebars.registerHelper({
		eq: (v1, v2) => v1 === v2,
		ne: (v1, v2) => v1 !== v2,
		lt: (v1, v2) => v1 < v2,
		gt: (v1, v2) => v1 > v2,
		lte: (v1, v2) => v1 <= v2,
		gte: (v1, v2) => v1 >= v2,
		and() {
			return Array.prototype.every.call(arguments, Boolean);
		},
		or() {
			return Array.prototype.slice.call(arguments, 0, -1).some(Boolean);
		}
	});

	function _render(view) {
		if (!view) {
			template.innerHTML = "Variable 'view' not found";
			return
		}
		var source = document.getElementById('template').innerHTML;
		var template = Handlebars.compile(source);
		var rendered = template(view);
		
		var el = document.createElement('div');
		el.id = '` + id.String() + `'
		el.innerHTML = rendered
		morphdom(document.getElementById('` + id.String() + `'), el);
	}
	var Embed = {
		render: _render,
		call: function(endpoint, request, callback, namespace) {
			if (!namespace) {
				namespace = "backend"
			}
			Micro.post(
				endpoint,
				namespace,
				request,
				function (data) {
					callback(data)
				}
			)
		},
		isLoggedIn: false,
		requireLogin: Micro.requireLogin,
		project: "` + project + `",
		user: {},
		isOwner: false
	}
	var _counter = 0;
	var _start = function() {
		_counter++
		if (_counter < 2) {
			return
		}
		` + jsFile + `
	}

	if (getCookie("micro_access")) {
	// if (false) {

		// triggering a refreshal of the token
		Embed.call("files/list", {
			project: "helloworld"
		}, function(dat) {
			Embed.call("auth/Auth/Inspect", {
				"options": {
					"namespace": "backend",
				},
				"token": getCookie("micro_access"),
			}, function(dat) {
				Embed.user = dat.account
				if (dat.account.id === "` + owner + `") {
					Embed.isOwner = true
				}
				Embed.isLoggedIn = true
				if (Embed.user.metadata) {
					Embed.user.name = Embed.user.metadata.username
				}
				_start();
			}, "micro")
		})
	} else {
		_counter++
	}

	document.addEventListener("DOMContentLoaded", function (event) {
		_start();
	})</script>` +
		`</body>
</html>`

	rsp.Header = make(map[string]*pb.Pair)
	rsp.Header["Content-Type"] = &pb.Pair{
		Key:    "Content-Type",
		Values: []string{"text/html", "charset=UTF-8"},
	}
	rsp.Body = rendered
	return nil
}

func extractScriptLink(html string) (script string, link string, rest string) {
	for _, line := range strings.Split(html, "\n") {
		switch {
		case strings.HasPrefix(line, "<script"):
			script += line
		case strings.HasPrefix(line, "<link"):
			link += line
		default:
			rest += line
		}
	}
	return
}
