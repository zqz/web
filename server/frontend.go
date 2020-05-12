package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var indexCache string
var pushStaticAssets []string

func (s Server) pushStatic(w http.ResponseWriter) {
	pusher, ok := w.(http.Pusher)
	if !ok {
		return
	}

	for _, a := range pushStaticAssets {
		pusher.Push(a, nil)
	}
}

func (s Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	if indexCache != "" {
		render.HTML(w, r, indexCache)
		s.pushStatic(w)
		return
	}

	var css, js []string
	var err error

	assets := assetFS()

	css, err = assets.AssetDir("build/static/css")
	if err != nil {
		panic(err)
	}
	js, err = assets.AssetDir("build/static/js")
	if err != nil {
		panic(err)
	}

	pushStaticAssets = make([]string, 0)
	for _, x := range css {
		pushStaticAssets = append(pushStaticAssets, "/assets/static/css/"+x)
	}

	for _, x := range js {
		pushStaticAssets = append(pushStaticAssets, "/assets/static/js/"+x)
	}

	var proto string

	if r.TLS != nil {
		proto = "https://"
	} else {
		proto = "http://"
	}

	var cdnRoot string
	cdnRoot = proto + r.Host + "/api/d/"

	tmplData := map[string]interface{}{
		"WSPath":  template.JSStr('/'),
		"ApiRoot": template.JSStr(fmt.Sprintf("%s%s/api", proto, r.Host)),
		"CdnRoot": template.JSStr(cdnRoot),

		"Assets": map[string]interface{}{
			"Js":  js,
			"Css": css,
		},
	}

	tmplContent := `
<!DOCTYPE HTML>
<html>
  <head>
    <meta http-equiv='content-type' content='text/html; charset=utf-8'>
    <title>zqz.ca</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">

    <link rel='shortcut icon' href='/assets/favicon.ico'/>
    {{- range .Assets.Css }}
    <link rel='stylesheet' media='screen' href='/assets/static/css/{{ . }}'/>
    {{- end }}
  </head>
  <body>
    <script type='text/javascript'>
      window.apiRoot = {{.ApiRoot}};
      window.cdnRoot = {{.CdnRoot}};
    </script>
  	<noscript>You need to enable JavaScript to run this app.</noscript><div id="root"></div>

    {{- range .Assets.Js }}
    <script type='text/javascript' src='/assets/static/js/{{.}}'></script>
    {{- end }}
  </body>
</html>`

	t := template.New("App Index Template")
	t, err = t.Parse(tmplContent)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	err = t.Execute(&output, tmplData)
	if err != nil {
		panic(err)
	}

	indexCache = output.String()
	render.HTML(w, r, indexCache)
	s.pushStatic(w)
}

func serveAssets(r chi.Router) {
	path := "/assets"
	root := assetFS()

	fs := http.StripPrefix(path, http.FileServer(root))

	// todo clean
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func fs(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		return
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
