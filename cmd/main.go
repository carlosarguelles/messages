package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/carlosarguelles/messages/internal"
	"github.com/carlosarguelles/messages/templates"
	"github.com/redis/go-redis/v9"
)

func main() {
	srv := http.NewServeMux()

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer client.Close()

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

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Fatal("error reading body")
			}
			message := r.Form.Get("message")
			fmt.Printf("Message received: %s", message)
		}
	})

	http.ListenAndServe(":8080", srv)
}
