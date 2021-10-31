package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

// https://pkg.go.dev/net/http#Handler
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("template", t.filename)))
	})

	err := t.template.Execute(w, nil)
	if err != nil {
		return
	}
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()

	log.Println("チャット開始！！")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
