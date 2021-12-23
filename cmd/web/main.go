package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// регаем роутер
	mux := http.NewServeMux()
	// роуты путь и обработчик
	mux.HandleFunc("/", home)
	mux.HandleFunc("/404", notFound)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// отдача статических файлов
	fileServer := http.FileServer(safeFileSystem{fs: http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Запуск веб-сервера на http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}

}

type safeFileSystem struct {
	fs http.FileSystem
}

func (nfs safeFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
