package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"os"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	err := t.template.Execute(w, data)
	if err != nil {
		return
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("環境変数が読み込めませんでした: ", err)
	}

	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	gomniauth.SetSecurityKey(os.Getenv("AUTH_SECURITY_KEY"))
	gomniauth.WithProviders(
		facebook.New("", "", "http://localhost:8080/auth/callback/facebook"),
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "http://localhost:8080/auth/callback/google"),
		github.New("", "", "http://localhost:8080/auth/callback/github"),
		)

	r := newRoom()
	//r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()

	log.Println("Webサーバーを開始します、ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
