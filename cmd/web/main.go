package main

import (
	"flag"
	"fmt"
	env "github.com/volodec/go-dot-env"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application struct {
	errLog  *log.Logger
	infoLog *log.Logger
}

func main() {
	// аналогия .env принимающая флаги из командной строки при запуске/сборке
	// пример для доступа по порту 3000: go run ./cmd/web -port="3000"
	port := flag.String("port", "4000", "Порт доступности приложения")
	flag.Parse()

	// использование самописа для получения значений из .env
	host := env.String("HOST", "localhost")
	addr := fmt.Sprintf("%s:%s", host, *port)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(findOrCreateFile("./ui/static/logs/errors.log"), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// логгер в файл ошибок http-сервера
	httpErrLog := log.New(findOrCreateFile("./ui/static/logs/httpErr.log"), "HTTP\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errLog:  errLog,
		infoLog: infoLog,
	}

	// Инициализация новой структуры http.Server для использования кастомного логгера ошибок
	srv := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: httpErrLog,
	}

	infoLog.Println(fmt.Sprintf("Запуск веб-сервера на http://%s:%s", host, *port))

	err := srv.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}

}

func findOrCreateFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}
	return file
}

// Обёртка над файловой системой.
type safeFileSystem struct {
	fs http.FileSystem
}

// Open реализует запрет на просмотр содержимого папки,
// с условием, что внутри нет index.html.
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
