package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// запрет на использование методов запросов отличных от POST
	if r.Method != http.MethodPost {
		// передаем в заголовке, какой метод разрешён
		w.Header().Set("Allow", http.MethodPost)

		// выдача кода состояния с описанием проблемы в теле ответа
		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	_, err := w.Write([]byte("Форма для создания новой заметки..."))
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, errRes := fmt.Fprintf(w, "Отображение заметки №%d...", id)
	if errRes != nil {
		app.errLog.Fatal(errRes)
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Так как в настройках роута последним и единственным символом является "/",
	// то текущий обработчик является обработчиком многоуровневого роута. Т.е.
	// сюда валится всё что не попадает в иные роуты текущего уровня.
	// Поэтому таким образом, мы реализуем 404 ошибку для несуществующих адресов.
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
