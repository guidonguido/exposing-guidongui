package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
)

type PostData struct {
	// Metadata fields
	Title   string
	Author  string
	URL     string
	Date    string
	Summary string

	// TODO: Check if a more efficient type than string is suitable for Content
	Content string

	ViewCount int
}

type PostReader interface {
	ReadMetadata(id string, data *PostData) error
	ReadContent(id string) ([]byte, error)
}

type HomepageData struct {
	Posts []PostData
}

type HomepageReader interface {
	ReadPosts() ([]PostData, error)
}

type FileReader struct{}

type IndexReader struct{}

// Read the Markdown content given a post ID
// Ignore metadata lines starting with {{ and ending with }}
func (fr FileReader) ReadContent(id string) ([]byte, error) {
	f, err := os.Open(fmt.Sprintf("web/posts/%s.md", id))
	if err != nil {
		slog.Error("Error opening post file", "error", err)
		return nil, err
	}
	defer f.Close()

	// Return bytes excluding metadata lines
	b := make([]byte, 0)

	reader := bufio.NewReader(f)
	// Read until newline (0xA) or EOF
	ln, err := reader.ReadBytes(0xA)
	for err == nil {
		// Check if the line is a metadata line (starts with "{{")
		if bytes.HasPrefix(ln, []byte(METADATA_PREFIX)) {
			// Read the next line
			ln, err = reader.ReadBytes(0xA)
			continue
		}

		b = append(b, ln...)
		// Continue reading until EOF
		ln, err = reader.ReadBytes(0xA)
	}

	return b, nil
}

func (ir IndexReader) ReadPosts() ([]PostData, error) {
	// Read all .md files in web/posts/
	dirEntries, err := os.ReadDir("web/posts/")
	if err != nil {
		slog.Error("Error reading posts directory", "error", err)
		return nil, err
	}

	posts := make([]PostData, 0, len(dirEntries))

	fr := FileReader{}

	for i := len(dirEntries) - 1; i >= 0; i-- {
		entry := dirEntries[i]
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if filename[len(filename)-3:] != ".md" {
			continue
		}
		id := filename[:len(filename)-3]

		var pd PostData

		// Read Metadata
		err := fr.ReadMetadata(id, &pd)
		if err != nil {
			slog.Error("Error reading post metadata", "id", id, "error", err)
			continue
		}

		pd.URL = "/posts/" + id

		posts = append(posts, pd)
	}

	return posts, nil
}
