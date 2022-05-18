package main

import (
	"errors"
	"fmt"
	"github.com/volodec/snippetbox/pkg/models"
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

	// TODO пока для теста
	title := "История про улитку"
	content := "Жила была улитка, пока не умерла"
	expires := "7"

	id, insetErr := app.snippets.Insert(title, content, expires)
	if insetErr != nil {
		app.serverError(w, insetErr)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}

		app.serverError(w, err)
		return
	}

	templateData := &templateData{OneEntry: snippet}

	app.render(w, r, "show.page.tmpl", templateData)
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

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	templateData := &templateData{ListEntries: snippets}

	app.render(w, r, "home.page.tmpl", templateData)
}
