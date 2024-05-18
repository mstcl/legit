package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/mstcl/legit/git"
)

func (d *deps) Write404(w http.ResponseWriter) {
	t := template.Must(template.ParseFS(Templates, filepath.Join("templates", "*")))
	w.WriteHeader(404)
	if err := t.ExecuteTemplate(w, "404", nil); err != nil {
		log.Printf("404 template: %s", err)
	}
}

func (d *deps) Write500(w http.ResponseWriter) {
	t := template.Must(template.ParseFS(Templates, filepath.Join("templates", "*")))
	w.WriteHeader(500)
	if err := t.ExecuteTemplate(w, "500", nil); err != nil {
		log.Printf("500 template: %s", err)
	}
}

func (d *deps) listFiles(files []git.NiceTree, data map[string]any, w http.ResponseWriter) {
	t := template.Must(template.ParseFS(Templates, filepath.Join("templates", "*")))

	data["files"] = files
	data["meta"] = d.c.Meta

	if err := t.ExecuteTemplate(w, "tree", data); err != nil {
		log.Println(err)
		return
	}
}

func (d *deps) showFile(content string, data map[string]any, w http.ResponseWriter) {
	t := template.Must(template.ParseFS(Templates, filepath.Join("templates", "*")))

	data["content"] = template.HTML(content)
	data["meta"] = d.c.Meta

	if err := t.ExecuteTemplate(w, "file", data); err != nil {
		log.Println(err)
		return
	}
}

func (d *deps) showRaw(content string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(content)); err != nil {
		log.Println(err)
	}
}
