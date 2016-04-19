package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
)

func handleTouch(w http.ResponseWriter, r *http.Request) {

	cid, err := r.Cookie(counterId)
	if err != nil {
		cookie := http.Cookie{
			Name:    counterId,
			Value:   uuid.NewV4().String(),
			Expires: time.Now().Add(365 * 24 * time.Hour),
		}
		cid = &cookie
		http.SetCookie(w, &cookie)
		log.Print(cookie)
	}

	v := &visitor{
		id:          cid.Value,
		lastTouched: time.Now().Unix(),
	}

	p.addVisitor(v)

	fmt.Fprint(w, strconv.Itoa(p.count()))

	log.Printf("touched: %s\n", cid.Value)

	return
}
