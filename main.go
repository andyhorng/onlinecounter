package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

const counterId = "counter_id"

var addr = flag.String("addr", ":8080", "http service address")

type visitor struct {
	id          string // uuid
	lastTouched int64  // unix timestamp
}

var testTempl = template.Must(template.ParseFiles("test.html"))

type pool struct {
	visitors map[string]*visitor
	lock     *sync.Mutex
}

var p = pool{
	visitors: make(map[string]*visitor),
	lock:     &sync.Mutex{},
}

func (p pool) addVisitor(v *visitor) int {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.visitors[v.id] = v

	connectionsLock.Lock()
	defer connectionsLock.Unlock()

	for c, _ := range connections {
		select {
		case c.ch <- true:
		default:
		}
	}

	return len(p.visitors)
}

func (p pool) count() int {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.visitors)
}

func (p pool) clean() {
	ticker := time.Tick(1 * time.Minute)

	for now := range ticker {
		log.Println("start clean")

		p.lock.Lock()
		expire := now.Add(-30 * time.Minute)

		for _, v := range p.visitors {
			log.Printf("%d %d\n", v.lastTouched, expire.Unix())
			if v.lastTouched < expire.Unix() {
				delete(p.visitors, v.id)

				log.Printf("delete: %s\n", v.id)
			}
		}
		p.lock.Unlock()
		log.Println("done clean")
	}

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

	go p.clean()

	http.ListenAndServe(*addr, nil)
}
