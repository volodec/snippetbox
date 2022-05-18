package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	env "github.com/volodec/go-dot-env"
	"github.com/volodec/snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errLog        *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// аналогия .env принимающая флаги из командной строки при запуске/сборке
	// пример для доступа по порту 3000: go run ./cmd/web -port="3000"
	port := flag.String("port", "4000", "Порт доступности приложения")

	// Определение нового флага из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "snippet_user:snippet_pass@/snippet_db?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	// использование самописа для получения значений из .env
	host := env.String("HOST", "localhost")
	addr := fmt.Sprintf("%s:%s", host, *port)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(findOrCreateFile("./ui/static/logs/errors.log"), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Fatal(err)
		return
	}

	app := &application{
		errLog:        errLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Инициализация новой структуры http.Server для использования кастомного логгера ошибок
	srv := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Println(fmt.Sprintf("Запуск веб-сервера на http://%s:%s", host, *port))

	listenErr := srv.ListenAndServe()
	if listenErr != nil {
		errLog.Fatal(listenErr)
	}

}

func findOrCreateFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}
	return file
}

// Функция openDB() обертывает sql.Open() и возвращает пул соединений sql.DB
// для заданной строки подключения (DSN).
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
