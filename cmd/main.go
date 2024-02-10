package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/carlosarguelles/messages/internal"
	"github.com/carlosarguelles/messages/templates"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{}

func main() {
	srv := http.NewServeMux()

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

	pool := internal.NewPool()

	go pool.Run()

	chatService := internal.NewChatService(client)

	srv.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	srv.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			templates.CreateChat().Render(ctx, w)
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal("error reading body")
			}
			name := r.Form.Get("name")
			chat, err := chatService.NewChat(ctx, name)
			if err != nil {
				log.Fatal("error creating chat")
			}
			http.Redirect(w, r, fmt.Sprintf("/chat?id=%s", chat.ID), http.StatusSeeOther)
		}
	})

	srv.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal("error reading body")
			}
			chatID := r.Form.Get("id")
			chat, err := chatService.GetChat(ctx, chatID)
			if err != nil {
				log.Fatal("error retrieving chat")
			}
			templates.Chat(*chat).Render(ctx, w)
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			content := r.Form.Get("message")
			username := r.Form.Get("username")
			chatID := r.Form.Get("chatID")
			message, err := chatService.NewMessage(r.Context(), chatID, username, content)
			if err != nil {
				templates.MessageForm(chatID, username, false).Render(r.Context(), w)
				return
			}
			pool.Broadcast <- *message
			templates.MessageForm(chatID, username, true).Render(r.Context(), w)
		}

		if r.Method == http.MethodPut {
			r.ParseForm()
			chatID := r.Form.Get("id")
			username := r.Form.Get("username")
			if username == "" || chatID == "" {
				return
			}
			templates.MessageForm(chatID, username, true).Render(r.Context(), w)
		}
	})

	srv.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		chatID := r.Form.Get("chatID")
		if chatID == "" {
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		client := internal.NewClient(pool, conn, chatID)
		pool.Register <- client
		go client.WritePump(r.Context())
	})

	http.ListenAndServe(":8080", srv)
}
