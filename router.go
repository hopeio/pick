package pick

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/liov/pick/utils"
)

const MethodAny = "*"

type muxEntry struct {
	preUrl string
	handle []*methodHandle
}

type groupMiddle struct {
	preUrl string
	handle []http.HandlerFunc
}

type methodHandle struct {
	method      string
	middleware  []http.HandlerFunc
	httpHandler http.Handler
	handle      reflect.Value
}

type Router struct {
	mu           sync.RWMutex
	route        map[string][]*methodHandle
	es           []muxEntry
	group        []groupMiddle
	NotFound     http.Handler
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
	hosts        bool
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

func (r *Router) ServeFiles(path string, root string) {

	fileServer := http.FileServer(http.Dir(root))

	r.es = appendSorted(r.es, muxEntry{
		path,
		[]*methodHandle{{
			http.MethodGet,
			nil,
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				req.URL.Path = req.URL.Path[len(path):]
				fileServer.ServeHTTP(w, req)
			}),
			reflect.Value{},
		},
		},
	})
}

func (r *Router) Use(middleware ...http.HandlerFunc) {
	for i, g := range r.group {
		if g.preUrl == "/" {
			r.group[i].handle = append(r.group[i].handle, middleware...)
			return
		}
	}
	r.group = append(r.group, groupMiddle{"/", middleware})
}

func (r *Router) Handle(method, path, describe string, handle ...http.HandlerFunc) {
	newMh := &methodHandle{method, handle[:len(handle)-1], handle[len(handle)-1], reflect.Value{}}
	if mh, ok := r.route[path]; ok {
		if h, _, _ := getHandle(method, mh); h != nil {
			panic("url：" + path + "已注册")
		}
	} else {
		r.route[path] = append(mh, newMh)
		if path[len(path)-1] == '/' {
			r.es = appendSorted(r.es, muxEntry{path, []*methodHandle{newMh}})
		}
	}

	if path[0] != '/' {
		r.hosts = true
	}
	fmt.Printf(" %s\t %s %s\t %s\n",
		utils.Green("API:"),
		utils.Yellow(utils.FormatLen(method, 6)),
		utils.Blue(utils.FormatLen(path, 50)), utils.Purple(describe))
}

func getHandle(method string, mhs []*methodHandle) (http.Handler, reflect.Value, []http.HandlerFunc) {
	if len(mhs) == 1 {
		if mhs[0].method == MethodAny || mhs[0].method == method {
			return mhs[0].httpHandler, mhs[0].handle, mhs[0].middleware
		}
		return nil, reflect.Value{}, nil
	}
	for _, mh := range mhs {
		if mh.method == method {
			return mh.httpHandler, mh.handle, mh.middleware
		}
	}
	return nil, reflect.Value{}, nil
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	for i, mh := range es {
		if mh.preUrl == e.preUrl {
			es[i].handle = append(es[i].handle, e.handle...)
			return es
		}
	}

	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].preUrl) < len(e.preUrl)
	})
	if i == n {
		return append(es, e)
	}

	es = append(es, muxEntry{})
	copy(es[i+1:], es[i:])
	es[i] = e
	return es
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}
	if mh, ok := r.route[req.URL.Path]; ok {
		h1, h2, middle := getHandle(req.Method, mh)
		if middle != nil {
			for _, f := range middle {
				f(w, req)
			}
		}
		for i := range r.group {
			if strings.HasPrefix(req.URL.Path, r.group[i].preUrl) {
				for _, f := range r.group[i].handle {
					f(w, req)
				}
			}
		}

		if h1 != nil {
			h1.ServeHTTP(w, req)
			return
		}
		if h2.IsValid() {
			commonHandler(w, req, h2)
			return
		}

	}
	for i := range r.es {
		if strings.HasPrefix(req.URL.Path, r.es[i].preUrl) {
			h1, _, middle := getHandle(req.Method, r.es[i].handle)
			if middle != nil {
				for _, f := range middle {
					f(w, req)
				}
			}
			h1.ServeHTTP(w, req)
			return
		}
	}

	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func appendGroupSorted(es []groupMiddle, e groupMiddle) []groupMiddle {
	for i, mh := range es {
		if mh.preUrl == e.preUrl {
			es[i].handle = append(es[i].handle, e.handle...)
			return es
		}
	}

	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].preUrl) < len(e.preUrl)
	})
	if i == n {
		return append(es, e)
	}

	es = append(es, groupMiddle{})
	copy(es[i+1:], es[i:])
	es[i] = e
	return es
}
