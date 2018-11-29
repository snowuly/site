package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"kob"
	"log"
	"net/http"
	"regexp"
	session "session-go"
	_ "session-memory"
)

type IndexData struct {
	IsLogin bool
	Name    string
	Login   string
}

var smng = session.NewManager("memory", "gsid", 24*3600)

func main() {
	defer db.Close()
	go smng.GC()

	var app kob.App

	indexTpl, err := template.ParseFiles("view/layout.tpl", "view/topbar.tpl", "view/index.tpl")
	checkErr(err)
	app.Get("/",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			store := ctx.Value("session").(session.Session)
			var uid int64
			var name, login string
			if value := store.Get("uid"); value != nil {
				uid = value.(int64)
			}
			if uid > 0 {
				err := db.QueryRow("select name, login from user where id = ?", uid).Scan(&name, &login)
				checkErr(err)
			}
			data := &IndexData{uid > 0, name, login}
			indexTpl.ExecuteTemplate(w, "index.tpl", data)
		},
	)

	app.Get("/register",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			registerTpl, err := template.ParseFiles("view/layout.tpl", "view/register.tpl")
			checkErr(err)
			registerTpl.ExecuteTemplate(w, "register.tpl", nil)
		},
	)

	app.Post("/register",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {

			login := r.FormValue("login")
			nickname := r.FormValue("nickname")
			pwd := r.FormValue("pwd")
			repwd := r.FormValue("repwd")
			loginReg := regexp.MustCompile("^[a-zA-Z0-9_\\-\\.]{2,20}$")

			var msg string
			switch {
			case !loginReg.MatchString(login):
				msg = "login name only can use letters, numbers, periods and length bettwen 2 an 20"
			case pwd == "":
				msg = "Password can't be empty."
			case pwd != repwd:
				msg = "Those passwords didn't match."
			}
			if msg != "" {
				w.Write([]byte(msg))
				return
			}
			stmt, err := db.Prepare("insert into user(login, name, pwd) values(?,?,?)")
			checkErr(err)
			if result, err := stmt.Exec(login, nickname, pwd); err == nil {
				id, _ := result.LastInsertId()
				store := ctx.Value("session").(session.Session)
				store.Set("uid", id)
				http.Redirect(w, r, "/", 302)
			} else {
				w.Write([]byte(fmt.Sprintf("Error: %s\n", err.Error())))
			}
		},
	)

	app.Get("/login",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			loginTpl, err := template.ParseFiles("view/layout.tpl", "view/login.tpl")
			checkErr(err)
			loginTpl.ExecuteTemplate(w, "login.tpl", nil)
		},
	)

	app.Post("/login",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			login := r.FormValue("login")
			pwd := r.FormValue("pwd")
			if login == "" || pwd == "" {
				w.Write([]byte("login or password cannot be empty"))
				return
			}
			var uid int64
			err := db.QueryRow("select id from user where login = ? and pwd = ?", login, pwd).Scan(&uid)
			if err != nil {
				if err == sql.ErrNoRows {
					w.Write([]byte("login or password wrong"))
					return
				}
				log.Fatal(err)
			}
			store := ctx.Value("session").(session.Session)
			store.Set("uid", uid)
			http.Redirect(w, r, "/", 302)
		},
	)

	app.Get("/logout",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			smng.SessionDestroy(w, r)
			http.Redirect(w, r, "/", 302)
		},
	)

	app.Get("/chat",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			store := ctx.Value("session").(session.Session)
			var uid int64
			if value := store.Get("uid"); value != nil {
				uid = value.(int64)
			}
			if uid == 0 {
				http.Redirect(w, r, "/login", 302)
				return
			}
			var name, login string

			err := db.QueryRow("select name, login from user where id = ?", uid).Scan(&name, &login)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			chatTpl, err := template.ParseFiles("view/layout.tpl", "view/topbar.tpl", "view/chat.tpl")
			checkErr(err)
			err = chatTpl.ExecuteTemplate(w, "chat.tpl", &IndexData{true, name, login})
		},
	)
	app.Post("/chat_msg",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			var uid int64
			store := ctx.Value("session").(session.Session)
			if value := store.Get("uid"); value != nil {
				uid = value.(int64)
			}
			if uid == 0 {
				http.Error(w, "no login", 403)
				return
			}
			msg := r.FormValue("msg")
			if msg == "" {
				http.Error(w, "empty msg", 400)
				return
			}
			chatmng.Send(uid, msg)
		},
	)

	app.Get("/chat_sse",
		sessionInit,
		func(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
			w.Header().Set("content-type", "text/event-stream; charset=utf-8")
			var uid int64
			store := ctx.Value("session").(session.Session)
			if value := store.Get("uid"); value != nil {
				uid = value.(int64)
			}
			if uid == 0 {
				http.Error(w, "not login", 403)
				return
			}
			var name, login string
			err := db.QueryRow("select name, login from user where id = ?", uid).Scan(&name, &login)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			receive, abort := chatmng.Enter(uid, name, login)

			w.Write([]byte("data:" + chatmng.GetUserList() + "\n\n"))
			w.(http.Flusher).Flush()
			for {
				select {
				case msg, ok := <-receive:
					if !ok {
						return
					}
					w.Write([]byte("data:" + msg + "\n\n"))
					w.(http.Flusher).Flush()
				case <-abort:
					return
				case <-r.Context().Done():
					chatmng.Leave(uid)
					return
				}

			}
		},
	)

	app.ListenTLS(":8080", "/Users/snow/self/openssl/raw.crt", "/Users/snow/self/openssl/raw.key")
	// app.Listen(":8080")
}

func sessionInit(ctx context.Context, w http.ResponseWriter, r *http.Request, next kob.NextFunc) {
	s := smng.SessionStart(w, r)
	w.Header().Set("content-type", "text/html; charset=utf-8")
	next(context.WithValue(ctx, "session", s))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
