package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}

	err := ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}
