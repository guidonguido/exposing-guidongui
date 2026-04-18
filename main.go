package main

import (
	"bytes"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"github.com/guidonguido/exposing-guidongui/internal"
)

func main() {
	blogMux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("web"))
	blogMux.Handle("GET /templates/scripts/", http.StripPrefix("/", fs))
	blogMux.Handle("GET /templates/styles/", http.StripPrefix("/", fs))
	blogMux.Handle("GET /posts/statics/", http.StripPrefix("/", fs))

	blogMux.HandleFunc("GET /", HomepageHandler(internal.IndexReader{}))
	blogMux.HandleFunc("GET /posts/{id}", PostHandler(internal.FileReader{}))

	slog.Info("Starting server on :8080")

	err := http.ListenAndServe(":8080", blogMux)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

// Handle Homepage index page
func HomepageHandler(rd internal.HomepageReader) http.HandlerFunc {

	tpl, err := template.ParseFiles("web/templates/home.gohtml")

	return func(w http.ResponseWriter, r *http.Request) {

		if err != nil {
			slog.Error("Template parsing error", "error", err)
			http.Error(w,
				"Error: unable to load template.",
				http.StatusInternalServerError)
			return
		}

		var hd internal.HomepageData
		hd.Posts, err = rd.ReadPosts()
		if err != nil {
			slog.Error("Error reading homepage posts", "error", err)
			http.Error(w,
				"Error: unable to read posts.",
				http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Render the template with HomepageData
		err = tpl.Execute(w, hd)
		if err != nil {
			slog.Error("Template execution error", "error", err)
			http.Error(w,
				"Error: unable to render template.",
				http.StatusInternalServerError)
			return
		}
	}
}

// Read a Post from whatever PostReader
func PostHandler(rd internal.PostReader) http.HandlerFunc {

	mdConverter := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)

	tpl, err := template.ParseFiles("web/templates/post.gohtml")

	return func(w http.ResponseWriter, r *http.Request) {

		// TRACE compute time
		// Start evaluation time
		renderTimeStart := time.Now()

		if err != nil {
			slog.Error("Template parsing error", "error", err)
			http.Error(w,
				"Error: unable to load template.",
				http.StatusInternalServerError)
			return
		}

		var id string = r.PathValue("id")

		var pd internal.PostData

		// 1. Read Post Metadata
		err := rd.ReadMetadata(id, &pd)
		if err != nil {
			slog.Error("Error reading post metadata", "id", id, "error", err)
			http.Error(w,
				"Error: unable to parse post metadata.",
				http.StatusInternalServerError)
			return
		}

		// 2. Read Post Content as Markdown
		mdPost, err := rd.ReadContent(id)
		if err != nil {
			slog.Error("Error reading post", "error", err)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// 3. Convert Markdown to HTML
		var buf bytes.Buffer
		if err := mdConverter.Convert(mdPost, &buf); err != nil {
			slog.Error("Markdown conversion error", "error", err)
			http.Error(w,
				"Error: unable to convert markdown to HTML.",
				http.StatusInternalServerError)
			return
		}
		pd.Content = buf.String()

		// TODO: 4. Get Post view count
		pd.ViewCount = 0

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// TRACE compute time
		// 5. Render the template with PostData
		// Empty buffer to fill with evaluated contet.
		// This is for computing time trace only. The http ResponseWriter buffer
		// is used for the http response instead
		var out bytes.Buffer
		err = tpl.Execute(&out, pd)
		if err != nil {
			slog.Error("Template execution error", "error", err)
			http.Error(w,
				"Error: unable to render template.",
				http.StatusInternalServerError)
			return
		}
		renderTime := time.Since(renderTimeStart)

		// 6. Write to output buffer
		nBytes, err := out.WriteTo(w)
		if err != nil {
			slog.Error("Write error", "error", err)
			http.Error(w,
				"Error: unable to write to output buffer.",
				http.StatusInternalServerError)
			return
		}

		slog.Info("Computing time info: ",
			"render time", renderTime,
			"# written bytes", nBytes)
	}
}
