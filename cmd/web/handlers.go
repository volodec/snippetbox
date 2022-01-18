package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// запрет на использование методов запросов отличных от POST
	if r.Method != http.MethodPost {
		// передаем в заголовке, какой метод разрешён
		w.Header().Set("Allow", http.MethodPost)

		str := fmt.Sprintf("%s-запрос запрещён. Разрешён только POST-запрос.", r.Method)
		// выдача кода состояния с описанием проблемы в теле ответа
		http.Error(w, str, http.StatusMethodNotAllowed)

		return
	}

	_, err := w.Write([]byte("Форма для создания новой заметки..."))
	if err != nil {
		app.errLog.Fatal(err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		redirectToNotFound(w, r)
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
		redirectToNotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Беда с чтением шаблонов...", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Беда с отрисовкой шаблонов...", http.StatusInternalServerError)
	}
}

func redirectToNotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404", http.StatusSeeOther)
}
