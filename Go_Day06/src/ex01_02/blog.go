package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	_ "github.com/lib/pq"
)

type Thought struct {
	ID      int
	Title   string
	Content string
	HTML    template.HTML
}

var (
	db        *sql.DB
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}).ParseGlob("templates/*.html"))
	adminUser    string
	adminPass    string
	postsPerPage = 3

	// ex02
	mu          sync.Mutex
	requests    = 0
	resetTicker = time.NewTicker(1 * time.Second)
	limit       = 100
)

func main() {
	loadCredentials()

	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if db, err = sql.Open("postgres", connStr); err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", rateLimitMiddleware(http.HandlerFunc(indexHandler)))
	http.Handle("/post", rateLimitMiddleware(http.HandlerFunc(postHandler)))
	http.Handle("/admin", rateLimitMiddleware(http.HandlerFunc(adminHandler)))

	fmt.Println("Server running at http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func loadCredentials() {
	file, err := os.ReadFile("admin_credentials.txt")
	if err != nil {
		log.Fatalf("ERROR: Something wrong with credentials file")
	}

	for _, line := range strings.Split(string(file), "\n") {
		if strings.HasPrefix(line, "username=") {
			adminUser = strings.TrimPrefix(line, "username=")
		}
		if strings.HasPrefix(line, "password=") {
			adminPass = strings.TrimPrefix(line, "password=")
		}
		if strings.HasPrefix(line, "db_name=") {
			os.Setenv("DB_NAME", strings.TrimPrefix(line, "db_name="))
		}
		if strings.HasPrefix(line, "db_user=") {
			os.Setenv("DB_USER", strings.TrimPrefix(line, "db_user="))
		}
		if strings.HasPrefix(line, "db_pass=") {
			os.Setenv("DB_PASS", strings.TrimPrefix(line, "db_pass="))
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * postsPerPage
	rows, err := db.Query(`
		SELECT id, title, content 
		FROM thoughts 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`,
		postsPerPage, offset)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var thoughts []Thought
	for rows.Next() {
		var t Thought
		if err := rows.Scan(&t.ID, &t.Title, &t.Content); err != nil {
			continue
		}
		mdParser := parser.NewWithExtensions(parser.CommonExtensions)
		renderer := html.NewRenderer(html.RendererOptions{})

		t.HTML = template.HTML(markdown.ToHTML([]byte(t.Content), mdParser, renderer))
		thoughts = append(thoughts, t)
	}

	templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Thoughts": thoughts,
		"Page":     page,
	})
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	idStr := path.Base(r.URL.Path)
	id, _ := strconv.Atoi(idStr)
	row := db.QueryRow(`
		SELECT title, content 
		FROM thoughts 
		WHERE id = $1`,
		id)

	var t Thought
	if err := row.Scan(&t.Title, &t.Content); err != nil {
		http.NotFound(w, r)
		return
	}
	mdParser := parser.NewWithExtensions(parser.CommonExtensions)
	renderer := html.NewRenderer(html.RendererOptions{})

	t.HTML = template.HTML(markdown.ToHTML([]byte(t.Content), mdParser, renderer))
	templates.ExecuteTemplate(w, "post.html", t)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")

		if user == adminUser && pass == adminPass {
			title := r.FormValue("title")
			content := r.FormValue("content")
			_, err := db.Exec(`
				INSERT INTO thoughts (title, content)
				VALUES ($1, $2)`,
				title, content)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			http.Redirect(w, r, "/", 302)
			return
		}
	}
	templates.ExecuteTemplate(w, "admin.html", nil)
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if requests >= limit {
			http.Error(w, "429 - Too Many Requests", http.StatusTooManyRequests)
			return
		}
		requests++
		next.ServeHTTP(w, r)
	})
}

func init() {
	go func() {
		for range resetTicker.C {
			mu.Lock()
			requests = 0
			mu.Unlock()
		}
	}()
}
