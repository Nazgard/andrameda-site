package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmplMain = template.Must(template.ParseFiles("template/index.html"))
var robots = template.Must(template.ParseFiles("static/robots.txt"))

func main() {
	// Подключаем статические файлы (включая favicon.ico)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Главная страница
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Начальная загрузка — просто отдаём HTML с пустой таблицей или с данными (можно пустую)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplMain.Execute(w, struct{ ShowResp bool }{ShowResp: r.URL.Query().Get("mode") == "resp"}); err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err := robots.ExecuteTemplate(w, "robots.txt", nil); err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			log.Println("Template execute error:", err)
		}
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
