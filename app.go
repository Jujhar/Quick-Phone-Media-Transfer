// ੴ ੨੦੧੭.੪.੨੬ ਸੁਭਮਸਤੁ
package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"fmt"
)

var templates = template.Must(template.ParseFiles("tmpl/backupnow.html"))
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

func downloadFromUrl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		display(w, "backupnow", nil)

	case "POST":

		// parse multipart
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get ref
		m := r.MultipartForm

		files := m.File["myfiles"]
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Write to destination
			dest, err := os.Create("./downloads/" + files[i].Filename)
			fmt.Println("Downloading", files[i].Filename, "to", "downloads/")
			defer dest.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// save uploaded file to dest
            n, err := io.Copy(dest, file)
            	if err != nil {
            		http.Error(w, err.Error(), http.StatusInternalServerError)
            		return
            	}
            	fmt.Println(n, "bytes downloaded.")

		}

		display(w, "backupnow", "Upload successful.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func main() {
    fmt.Println("running at UrIPAddress/backupnow..")
	http.HandleFunc("/backupnow", downloadFromUrl)
	http.ListenAndServe(":80", nil)
}
