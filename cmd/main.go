package main

import (
	"fmt"
	"net/http"

	"github.com/carlosarguelles/messages/internal"
	"github.com/carlosarguelles/messages/templates"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{}

func main() {
	srv := http.NewServeMux()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

	pool := internal.NewPool()

	defer pool.Close()

	go pool.Run()

	chatService := internal.NewChatService(client)

	srv.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	srv.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			templates.CreateChat().Render(r.Context(), w)
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			name := r.Form.Get("name")
			if name == "" {
				return
			}
			chat, err := chatService.NewChat(r.Context(), name)
			if err != nil {
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/chat?id=%s", chat.ID), http.StatusSeeOther)
		}
	})

	srv.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()
			chatID := r.Form.Get("id")
			chat, err := chatService.GetChat(r.Context(), chatID)
			if err != nil {
				http.Redirect(w, r, "/new", http.StatusSeeOther)
				return
			}
			username := r.Form.Get("username")
			messages, err := chatService.GetChatMessages(r.Context(), chatID)
			if err != nil {
				return
			}
			templates.Chat(*chat, username, messages).Render(r.Context(), w)
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
