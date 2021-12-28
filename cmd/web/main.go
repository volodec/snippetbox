package main

import (
	"flag"
	"fmt"
	env "github.com/volodec/go-dot-env"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// аналогия .env принимающая флаги из командной строки при запуске/сборке
	// пример для доступа по порту 3000: go run ./cmd/web -port="3000"
	port := flag.String("port", "4000", "Порт доступности приложения")
	flag.Parse()

	// использование самописа для получения значений из .env
	host := env.String("HOST", "localhost")

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

	log.Println(fmt.Sprintf("Запуск веб-сервера на http://%s:%s", host, *port))

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, *port), mux)
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
