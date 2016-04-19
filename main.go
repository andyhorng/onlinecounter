package main

import (
	"flag"
	"html/template"
	"net/http"
	"sync"
)

const counterId = "counter_id"

var addr = flag.String("addr", ":8080", "http service address")

type visitor struct {
	id          string // uuid
	lastTouched int64  // unix timestamp
}

var testTempl = template.Must(template.ParseFiles("test.html"))

type pool struct {
	visitors      map[string]*visitor
	lock          *sync.Mutex
	eventChannels []chan *visitor
}

var p = pool{
	visitors: make(map[string]*visitor),
	lock:     &sync.Mutex{},
}

func (p pool) addVisitor(v *visitor) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.visitors[v.id] = v

	connectionsLock.Lock()
	defer connectionsLock.Unlock()

	for c, _ := range connections {
		c.ch <- true
	}

}

func (p pool) count() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.visitors)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	testTempl.Execute(w, r.Host)
}

func main() {
	flag.Parse()

	http.HandleFunc("/touch", handleTouch)
	http.HandleFunc("/listen", handleListen)
	http.HandleFunc("/test", handleTest)

	http.ListenAndServe(*addr, nil)
}
