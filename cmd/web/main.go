package main

import (
	"log"
	"net/http"
)

func main() {
	// регаем роутер
	mux := http.NewServeMux()
	// роуты путь и обработчик
	mux.HandleFunc("/", home)
	mux.HandleFunc("/404", notFound)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Запуск веб-сервера на http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}

}
