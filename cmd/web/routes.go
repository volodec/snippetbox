package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// регаем роутер
	mux := http.NewServeMux()
	// роуты путь и обработчик
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// отдача статических файлов
	fileServer := http.FileServer(safeFileSystem{fs: http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
