package main

import (
	"html/template"
	"net/http"
	"os"
	"regexp"
)

// PageはWikiの各ページを表します
type Page struct {
	Title string
	Body  []byte
}

// ページ内容をファイルに保存する
func (p *Page) save() error {
	filename := "text_dir/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// ファイルからページ内容を読み込む
func loadPage(title string) (*Page, error) {
	filename := "text_dir/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// テンプレートのパース（edit.html, view.html）
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// 指定のテンプレートにデータを適用してレンダリングする
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// htmlにGoのpデータを渡している
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// URLパスの正規表現（ページ名は英数字のみ）
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// ハンドラーをラップし、URLからタイトル部分を抽出する
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// ページの閲覧
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// ページの編集
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// ページの保存
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
