package app

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type engineX struct {
	ctx context.Context

	localStorage   BrowserStorage
	sessionStorage BrowserStorage

	initBrowserOnce sync.Once
	browser         browser

	routes         *router
	internalURLs   []string
	resolveURL     func(string) string
	originPage     requestPage
	lastVisitedURL *url.URL

	nodes   nodeManager
	updates updateManager
	newBody func() HTMLBody
	body    UI

	dispatches chan func()
	defers     chan func()
	goroutines sync.WaitGroup

	asynchronousActionHandlers map[string]ActionHandler
	actions                    actionManager
	states                     stateManager
}

func newEngineX(ctx context.Context, routes *router, resolveURL func(string) string, origin *url.URL, newBody func() HTMLBody, actionHandlers map[string]ActionHandler) *engineX {
	var localStorage BrowserStorage
	var sessionStorage BrowserStorage
	if IsServer {
		localStorage = newMemoryStorage()
		sessionStorage = newMemoryStorage()
	} else {
		localStorage = newJSStorage("localStorage")
		sessionStorage = newJSStorage("sessionStorage")
	}

	if resolveURL == nil {
		resolveURL = func(v string) string { return v }
	}

	return &engineX{
		ctx:        ctx,
		routes:     routes,
		resolveURL: resolveURL,
		originPage: requestPage{
			url:                   origin,
			resolveStaticResource: resolveURL,
		},
		localStorage:               localStorage,
		lastVisitedURL:             &url.URL{},
		sessionStorage:             sessionStorage,
		newBody:                    newBody,
		nodes:                      nodeManager{},
		dispatches:                 make(chan func(), 4096),
		defers:                     make(chan func(), 4096),
		asynchronousActionHandlers: actionHandlers,
	}
}

func (e *engineX) baseContext() nodeContext {
	return nodeContext{
		Context:               e.ctx,
		resolveURL:            e.resolveURL,
		appUpdatable:          e.browser.AppUpdatable,
		page:                  e.page,
		navigate:              e.Navigate,
		localStorage:          e.localStorage,
		sessionStorage:        e.sessionStorage,
		dispatch:              e.dispatch,
		defere:                e.defere,
		async:                 e.async,
		addComponentUpdate:    e.updates.Add,
		removeComponentUpdate: e.updates.Done,
		handleAction:          e.actions.Handle,
		postAction:            e.actions.Post,
	}
}

// Navigate directs the engine to the specified URL destination, which might be
// an internal page within the app, an external link outside the app, or a
// mailto link. If the 'updateHistory' flag is true, the destination is added to
// the browser's history.
func (e *engineX) Navigate(destination *url.URL, updateHistory bool) {
	e.initBrowserOnce.Do(e.initBrowser)

	switch {
	case e.internalURL(destination),
		e.mailTo(destination):
		Window().Get("location").Set("href", destination.String())
		return

	case e.externalNavigation(destination):
		Window().Call("open", destination.String())
		return
	}

	if destination.String() == e.lastVisitedURL.String() {
		return
	}
	defer func() {
		if updateHistory {
			Window().addHistory(destination)
		}
		e.lastVisitedURL = destination

		e.nodes.NotifyComponentEvent(e.baseContext(), e.body, nav{})
		if destination.Fragment != "" {
			e.defere(func() {
				Window().ScrollToID(destination.Fragment)
			})
		}
	}()

	if destination.Path == e.lastVisitedURL.Path &&
		destination.Fragment != e.lastVisitedURL.Fragment {
		return
	}

	path := strings.TrimPrefix(destination.Path, Getenv("GOAPP_ROOT_PREFIX"))
	if path == "" {
		path = "/"
	}
	root, ok := e.routes.createComponent(path)
	if !ok {
		root = &notFound{}
	}
	e.load(root)
}

func (e *engineX) initBrowser() {
	if IsServer {
		return
	}
	e.browser.HandleEvents(e.baseContext(), e.notifyComponentEvent)
}

func (e *engineX) notifyComponentEvent(event any) {
	e.nodes.NotifyComponentEvent(e.baseContext(), e.body, event)
}

func (e *engineX) externalNavigation(v *url.URL) bool {
	return v.Host != e.originPage.URL().Host
}

func (e *engineX) mailTo(v *url.URL) bool {
	return v.Scheme == "mailto"
}

func (e *engineX) internalURL(v *url.URL) bool {
	if e.internalURLs == nil {
		json.Unmarshal([]byte(Getenv("GOAPP_INTERNAL_URLS")), &e.internalURLs)
	}

	url := v.String()
	for _, u := range e.internalURLs {
		if strings.HasPrefix(url, u) {
			return true
		}
	}
	return false
}

func (e *engineX) page() Page {
	if IsClient {
		return browserPage{resolveStaticResource: e.resolveURL}
	}
	return &e.originPage
}

func (e *engineX) load(v Composer) {
	if e.body == nil {
		body, err := e.nodes.Mount(e.baseContext(), 0, e.newBody().privateBody(v))
		if err != nil {
			panic(errors.New("mounting root failed").Wrap(err))
		}

		for action, handler := range e.asynchronousActionHandlers {
			e.actions.Handle(action, body, true, handler)
		}

		e.body = body
		return
	}

	body, err := e.nodes.Update(e.baseContext(), e.body, e.newBody().privateBody(v))
	if err != nil {
		panic(errors.New("updating root failed").Wrap(err))
	}
	e.body = body
}

func (e *engineX) dispatch(v func()) {
	e.dispatches <- v
}

func (e *engineX) defere(v func()) {
	e.defers <- v
}

func (e *engineX) async(v func()) {
	e.goroutines.Add(1)
	go func() {
		v()
		e.goroutines.Done()
	}()
}

// Start initiates the main event loop of the engine at the specified framerate.
// The loop efficiently manages dispatches, component updates, and deferred
// actions.
func (e *engineX) Start(framerate int) {
	if framerate <= 0 {
		framerate = 30
	}

	iddleFrameDuration := time.Hour
	activeFrameDuration := time.Second / time.Duration(framerate)
	currentFrameDuration := time.Nanosecond
	frames := time.NewTicker(currentFrameDuration)
	defer frames.Stop()

	for {
		select {
		case dispatch := <-e.dispatches:
			if currentFrameDuration != activeFrameDuration {
				frames.Reset(activeFrameDuration)
				currentFrameDuration = activeFrameDuration
			}
			dispatch()

		case <-frames.C:
			e.updates.ForEach(func(c Composer) {
				if _, err := e.nodes.UpdateComponentRoot(e.baseContext(), c); err != nil {
					panic(errors.New("updating component failed").Wrap(err))
				}
				e.updates.Done(c)
			})
			e.executeDefers()
			e.actions.Cleanup()
			e.states.Cleanup()

			frames.Reset(iddleFrameDuration)
			currentFrameDuration = iddleFrameDuration

		case <-e.ctx.Done():
			return
		}
	}
}

func (e *engineX) executeDefers() {
	for {
		select {
		case defere := <-e.defers:
			defere()

		default:
			return
		}
	}
}
