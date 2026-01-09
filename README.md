# exposing-guidongui

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-Apache%202.0-green)
![Self-Hosted](https://img.shields.io/badge/Self--Hosted-Home%20Lab-orange)

## Table of Contents

- [Motivation](#motivation)
- [Guidonguido Go Blog](#guidonguido-go-blog)
  - [net/http Standard Library](#nethttp-standard-library)
  - [Go Templates](#go-templates)
  - [Goldmark Markdown Parsing](#goldmark-markdown-parsing)
  - [Custom Metadata Parser](#custom-metadata-parser)
  - [Asciinema Integration](#asciinema-integration)

---

## Motivation

This repository contains the tech stack of my personal blog, reachable at exposing.guidongui.com.
As the subdomain name suggests, I'm not seeking fame, just a public place to talk about my personal projects. Public exposure is a way to do better what I do every day and above all a way to finish the small projects I start and then very often leave in POC state.

I've always wanted a space where I could pour out my ideas, but for deontological reasons, I've always refused to use CMS like Hugo to build one. The day I would do it, it would be built from scratch by me.
The goal is not to look pretty or smart to the world, but to learn by experimenting. I preferred simplicity over flashiness, custom Go libraries over frameworks.

The blog is completely self-developed, self-hosted in my Home Media Lab, which I hope I'll have the chance to talk about in future posts.

<details>
<summary>How this documentation was generated</summary>

Documentation for the blog is in development, but anticipating the moment when my creative development will be replaced by an LLM, I recommend cloning this repository and asking, including the codebase in the context, what the code does, diving deeper into:
- Usage of the Go net/http standard library
- Usage of Go templates
- External libraries used: Goldmark
- Custom metadata parsers in ./internal/metadata_parser.go
- Asciinema integration, including self-hosted server

**Example prompt:**

> Analyze the exposing-guidongui codebase and produce a concise and technical README.md documentation file, paying particular attention to the following points:
> - Usage of the Go net/http standard library
> - Usage of Go templates to render blog pages
> - External libraries: Goldmark usage for Markdown parsing
> - Custom metadata parsers in ./internal/metadata_parser.go
> - Asciinema integration, including self-hosted server

</details>

---

## Guidonguido Go Blog

### net/http Standard Library

The server uses Go's standard `net/http` package without external routing frameworks.

**Key patterns:**
- Uses Go 1.22+ method-based routing (`GET /path`)
- Path parameters extracted via `r.PathValue("id")`
- Handlers return `http.HandlerFunc` closures with dependency injection for readers
- `http.FileServer` serves static CSS and JS from `web/templates/`

### Go Templates

Templates use the `text/template` package with `.gohtml` files.

**Homepage template** (`home.gohtml`) iterates over posts:
```html
{{range .Posts}}
<article>
    <h3><a href="{{.URL}}">{{.Title}}</a></h3>
    <address>{{.Date}}</address>
    <pre>{{.Summary}}</pre>
</article>
{{end}}
```

**Post template** (`post.gohtml`) renders single post data:
```html
<title>{{.Title}}</title>
<pre>{{.Content}}</pre>
```

**Data structures passed to templates:**
```go
type HomepageData struct {
    Posts []PostData
}

type PostData struct {
    Title, Author, URL, Date, Summary, Content string
    ViewCount int
}
```

Templates are parsed once at handler creation and executed per-request via `tpl.Execute(w, data)`.

### Goldmark Markdown Parsing

[Goldmark](https://github.com/yuin/goldmark) converts Markdown posts to HTML.

```go
mdConverter := goldmark.New(
    goldmark.WithExtensions(extension.GFM),  // GitHub Flavored Markdown
)

var buf bytes.Buffer
mdConverter.Convert(mdPost, &buf)
pd.Content = buf.String()
```

The GFM extension enables:
- Tables
- Strikethrough (`~~text~~`)
- Task lists (`- [x] done`)
- Autolinks

### Custom Metadata Parser

Posts use a custom metadata format at the file start:

```markdown
{{Title: My Post Title}}
{{Author: Guidongui}}
{{Date: 04/01/2026}}
{{Summary: Brief description...}}
```

### Asciinema Integration

The homepage displays a terminal recording via a self-hosted Asciinema server.

**Frontend integration** (`home.gohtml`):
```javascript
AsciinemaPlayer.create(
    'https://asciinema.guidongui.com/a/34.cast',
    document.getElementById('asciinema-player'),
    { autoPlay: true, loop: true, theme: 'white', ... }
);
```

**Self-hosted server** (`deploy/asciinema/deployment.yaml`):
- Runs the official `ghcr.io/asciinema/asciinema-server` image
- PostgreSQL database in a sidecar container
- Exposed at `asciinema.guidongui.com` via Traefik IngressRoute with TLS
- Configuration: signups disabled, uploads require auth, public recordings

This allows uploading terminal recordings from the home lab and embedding them in blog posts without relying on third-party hosting.