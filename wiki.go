package main

import (
	"net/http"
	"regexp"
	"text/template"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// PageはWikiの各ページを表します
type Page struct {
	Title string
	Body  []byte
}

// DBのエンティティ
type HtmlEntity struct {
	ID    int    `gorm:"primaryKey"`
	Title string `gorm:"unique"`
	Body  string
}

type WikiDB struct {
	DB     *gorm.DB
	dbName string
}

func NewWikiDB(db *gorm.DB, dbName string) *WikiDB {
	return &WikiDB{DB: db, dbName: dbName}
}

// func (HtmlEntity) TableName() string {
// 	return ""
// }

// ページをDBに保存（新規 or 更新）
func (w *WikiDB) SavePage(title string, body []byte) error {
	var entity HtmlEntity
	result := w.DB.Where("title = ?", title).First(&entity)
	if result.Error == nil {
		entity.Body = string(body)
		return w.DB.Save(&entity).Error
	}
	entity = HtmlEntity{Title: title, Body: string(body)}
	return w.DB.Create(&entity).Error
}

// ページをDBから取得
func (w *WikiDB) LoadPage(title string) (*Page, error) {
	var entity HtmlEntity
	if err := w.DB.Where("title = ?", title).First(&entity).Error; err != nil {
		return nil, err
	}
	return &Page{Title: entity.Title, Body: []byte(entity.Body)}, nil
}

// テンプレート
func getTemplates() *template.Template {
	return template.Must(template.ParseFiles("html/edit.html", "html/view.html"))
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	templates := getTemplates()
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// URLパスの正規表現（ページ名は英数字のみ）
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// ハンドラーをラップし、URLからタイトル部分を抽出する
func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	// MySQL接続例（ユーザー名・パスワード・DB名は適宜変更）
	dsn := "test_user:nakanishi@tcp(localhost:3306)/test_db"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&HtmlEntity{})
	wikiDB := NewWikiDB(db, "test_db")

	http.HandleFunc("/view/", MakeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		p, err := wikiDB.LoadPage(title)
		if err != nil {
			http.Redirect(w, r, "/edit/"+title, http.StatusFound)
			return
		}
		renderTemplate(w, "view", p)
	}))
	http.HandleFunc("/edit/", MakeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		p, err := wikiDB.LoadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}
		renderTemplate(w, "edit", p)
	}))
	http.HandleFunc("/save/", MakeHandler(func(w http.ResponseWriter, r *http.Request, title string) {
		body := r.FormValue("body")
		err := wikiDB.SavePage(title, []byte(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	}))
	http.ListenAndServe(":8080", nil)
}
