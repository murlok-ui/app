// +build !wasm

package app

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

const (
	defaultThemeColor = "#2d2c2c"
)

// Handler is an HTTP handler that serves an HTML page that loads a Go wasm app
// and its resources.
type Handler struct {
	// The page authors.
	Author string

	// A placeholder background color for the application page to display before
	// its stylesheets are loaded.
	//
	// DEFAULT: #2d2c2c.
	BackgroundColor string

	// The page description.
	Description string

	// The icon that is used for the PWA, favicon, loading and default not
	// found component.
	//
	// DEFAULT: app default and large icons.
	Icon Icon

	// The page keywords.
	Keywords []string

	// The text displayed while loading a page.
	LoadingLabel string

	// The name of the web application as it is usually displayed to the user.
	Name string

	// The paths or urls of the JavaScript files to use with the page.
	//
	// Paths are relative to the web directory.
	Scripts []string

	// The name of the web application displayed to the user when there is not
	// enough space to display Name.
	ShortName string

	// The paths or urls of the CSS files to use with the page.
	//
	// Paths are relative to the web directory.
	Styles []string

	// The theme color for the application. This affects how the OS displays the
	// app (e.g., PWA title var or Android's task switcher).
	//
	// DEFAULT: #2d2c2c.
	ThemeColor string

	// The page title.
	Title string

	// The version number. This is used in order to update the pwa application
	// in the browser. It must be set when deployed on a scalable system in
	// order to prevent constant updates.
	//
	// Default: Auto-generated in order to trigger pwa update on a local
	// development sytem.
	Version string

	// The path of the directory where static resources like wasm program,
	// images, styles or scripts are located.
	//
	// DEFAULT: "web".
	Web string

	// The path the wasm program.
	//
	// Path is relative to the web directory.
	//
	// DEFAULT: "app.wasm".
	Wasm string

	once         sync.Once
	etag         string
	page         bytes.Buffer
	manifestJSON bytes.Buffer
	viJS         bytes.Buffer
	viWorkerJS   bytes.Buffer
	wasmExecJS   []byte
	viCSS        []byte
}

func (h *Handler) init() {
	h.initVersion()
	h.initWebDir()
	h.initStyles()
	h.initScripts()
	h.initIcon()
	h.initPWA()
	h.initPage()
	h.initWasmJS()
	h.initViJS()
	h.initWorkerJS()
	h.initManifestJSON()
	h.initScripts()
	h.initViCSS()
}

func (h *Handler) initVersion() {
	if h.Version == "" {
		t := time.Now().UTC().String()
		h.Version = fmt.Sprintf(`%x`, sha1.Sum([]byte(t)))
	}
	h.etag = `"` + h.Version + `"`
}

func (h *Handler) initWebDir() {
	if h.Web == "" {
		h.Web = "web"
	}
	h.Web = strings.TrimPrefix(h.Web, ".")
	h.Web = strings.TrimPrefix(h.Web, "/")
	h.Web = strings.TrimSuffix(h.Web, "/")
}

func (h *Handler) initStyles() {
	for i, path := range h.Styles {
		h.Styles[i] = h.staticResource(path)
	}
}

func (h *Handler) initScripts() {
	for i, path := range h.Scripts {
		h.Scripts[i] = h.staticResource(path)
	}
}

func (h *Handler) initIcon() {
	if h.Icon.Default == "" {
		h.Icon.Default = "https://storage.googleapis.com/murlok-github/icon-192.png"
		h.Icon.Large = "https://storage.googleapis.com/murlok-github/icon-512.png"
	}

	if h.Icon.AppleTouch == "" {
		h.Icon.AppleTouch = h.Icon.Default
	}

	h.Icon.Default = h.staticResource(h.Icon.Default)
	h.Icon.Large = h.staticResource(h.Icon.Large)
	h.Icon.AppleTouch = h.staticResource(h.Icon.AppleTouch)
}

func (h *Handler) initPWA() {
	if h.Name == "" && h.ShortName == "" {
		h.Name = "App PWA"
	}
	if h.ShortName == "" {
		h.ShortName = h.Name
	}
	if h.Name == "" {
		h.Name = h.ShortName
	}

	if h.BackgroundColor == "" {
		h.BackgroundColor = defaultThemeColor
	}
	if h.ThemeColor == "" {
		h.ThemeColor = defaultThemeColor
	}

	if h.LoadingLabel == "" {
		h.LoadingLabel = "Loading"
	}
}

func (h *Handler) initPage() {
	h.page.WriteString("<!DOCTYPE html>\n")
	Html().
		Body(
			Head().
				Body(
					Meta().Charset("UTF-8"),
					Meta().
						HTTPEquiv("Content-Type").
						Content("text/html; charset=utf-8"),
					Meta().
						Name("author").
						Content(h.Author),
					Meta().
						Name("description").
						Content(h.Description),
					Meta().
						Name("keywords").
						Content(strings.Join(h.Keywords, ", ")),
					Meta().
						Name("viewport").
						Content("width=device-width, initial-scale=1, maximum-scale=1, user-scalable=0"),
					Title().
						Body(
							Text(h.Title),
						),
					Link().
						Rel("icon").
						Type("image/png").
						Href(h.Icon.Default),
					Link().
						Rel("apple-touch-icon").
						Href(h.Icon.AppleTouch),
					Link().
						Rel("manifest").
						Href("/manifest.json"),
					Link().
						Type("text/css").
						Rel("stylesheet").
						Href("/app.css"),
					Range(h.Styles).
						Slice(func(i int) Node {
							return Link().
								Type("text/css").
								Rel("stylesheet").
								Href(h.Styles[i])
						}),
					Script().Src("/wasm_exec.js"),
					Script().Src("/app.js"),
					Range(h.Scripts).
						Slice(func(i int) Node {
							return Script().
								Src(h.Scripts[i])
						}),
				),
			Body().
				Body(
					Div().
						Class("app-wasm-layout").
						Body(
							Img().
								ID("app-wasm-loader-icon").
								Class("app-wasm-icon app-spin").
								Src(h.Icon.Default),
							P().
								ID("app-wasm-loader-label").
								Class("app-wasm-label").
								Body(Text(h.LoadingLabel)),
						),
					Div().ID("app-context-menu"),
				),
		).
		html(&h.page)
}

func (h *Handler) initWasmJS() {
	h.wasmExecJS = []byte(wasmExecJS)
}

func (h *Handler) initViJS() {
	if h.Wasm == "" {
		h.Wasm = "app.wasm"
	}
	h.Wasm = h.staticResource(h.Wasm)

	if err := template.
		Must(template.New("app.js").Parse(viJS)).
		Execute(&h.viJS, struct {
			Wasm string
		}{
			Wasm: h.Wasm,
		}); err != nil {
		panic(err)
	}
}

func (h *Handler) initWorkerJS() {
	cacheableResources := make(map[string]struct{})
	for _, s := range h.Scripts {
		cacheableResources[s] = struct{}{}
	}
	for _, s := range h.Styles {
		cacheableResources[s] = struct{}{}
	}
	cacheableResources["/wasm_exec.js"] = struct{}{}
	cacheableResources["/app.js"] = struct{}{}
	cacheableResources["/manifest.json"] = struct{}{}
	cacheableResources[h.Wasm] = struct{}{}
	cacheableResources[h.Icon.Default] = struct{}{}
	cacheableResources[h.Icon.Large] = struct{}{}
	cacheableResources[h.Icon.AppleTouch] = struct{}{}
	cacheableResources["/"] = struct{}{}

	if err := template.
		Must(template.New("app-worker.js").Parse(viWorkerJS)).
		Execute(&h.viWorkerJS, struct {
			Version          string
			ResourcesToCache map[string]struct{}
		}{
			Version:          h.Version,
			ResourcesToCache: cacheableResources,
		}); err != nil {
		panic(err)
	}
}

func (h *Handler) initManifestJSON() {
	if err := template.
		Must(template.New("manifest.json").Parse(manifestJSON)).
		Execute(&h.manifestJSON, struct {
			ShortName       string
			Name            string
			DefaultIcon     string
			LargeIcon       string
			BackgroundColor string
			ThemeColor      string
		}{
			ShortName:       h.ShortName,
			Name:            h.Name,
			DefaultIcon:     h.Icon.Default,
			LargeIcon:       h.Icon.Large,
			BackgroundColor: h.BackgroundColor,
			ThemeColor:      h.ThemeColor,
		}); err != nil {
		panic(err)
	}
}

func (h *Handler) initViCSS() {
	h.viCSS = []byte(viCSS)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.once.Do(h.init)

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", h.etag)

	etag := r.Header.Get("If-None-Match")
	if etag == h.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	switch r.URL.Path {
	case "/wasm_exec.js":
		h.serveWasmExecJS(w, r)
		return

	case "/app.js":
		h.serveViJS(w, r)
		return

	case "/app-worker.js":
		h.serveViWorkerJS(w, r)
		return

	case "/manifest.json":
		h.serveManifestJSON(w, r)
		return

	case "/app.css":
		h.serveViCSS(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/"+h.Web) {
		filename := strings.TrimPrefix(r.URL.Path, "/")
		filename = normalizeFilePath(filename)

		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			http.ServeFile(w, r, filename)
			return
		}
	}

	h.servePage(w, r)
}

func (h *Handler) servePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.page.Len()))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(h.page.Bytes())
}

func (h *Handler) serveWasmExecJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(len(h.wasmExecJS)))
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(h.wasmExecJS)
}

func (h *Handler) serveViJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.viJS.Len()))
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(h.viJS.Bytes())
}

func (h *Handler) serveViWorkerJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.viWorkerJS.Len()))
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(h.viWorkerJS.Bytes())
}

func (h *Handler) serveManifestJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(h.manifestJSON.Len()))
	w.Header().Set("Content-Type", "application/manifest+json")
	w.WriteHeader(http.StatusOK)
	w.Write(h.manifestJSON.Bytes())
}

func (h *Handler) serveViCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", strconv.Itoa(len(h.viCSS)))
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write(h.viCSS)
}

func (h *Handler) staticResource(path string) string {
	u, _ := url.Parse(path)
	if u.Scheme != "" {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "/" + normalizePath(h.Web) + path
}

func normalizePath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, `\`, "/")
	}
	return path
}

func normalizeFilePath(path string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(path, "/", `\`)
	}
	return path
}

// Icon describes a square image that is used in various places such as
// application icon, favicon or loading icon.
type Icon struct {
	// The path or url to a square image/png file. It must have a side of 192px.
	//
	// Paths are relative to the web directory.
	Default string

	// The path or url to larger square image/png file. It must have a side of
	// 512px.
	//
	// Path is relative to the web directory.
	Large string

	// The path or url to a square image/png file that is used for IOS/IPadOS
	// home screen icon. It must have a side of 192px.
	//
	// Path is relative to the web directory.
	//
	// DEFAULT: Icon.Default
	AppleTouch string
}
