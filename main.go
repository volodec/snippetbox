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

func notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func createSnippet(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Форма для создания новой заметки..."))
	if err != nil {
		log.Fatal(err)
	}
}

func showSnippet(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Отображение заметки..."))
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// Так как в настройках роута последним и единственным символом является "/",
	// то текущий обработчик является обработчиком многоуровневого роута. Т.е.
	// сюда валится всё что не попадает в иные роуты текущего уровня.
	// Поэтому таким образом, мы реализуем 404 ошибку для несуществующих адресов.
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	_, err := w.Write([]byte("Привет из SnippetBox"))
	if err != nil {
		log.Fatal(err)
	}
}
