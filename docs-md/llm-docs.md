---
created_on: 2024-12-26
template: "llm-docs.html"
---

# LLM-Friendly Documentation

We are moving in a fast paced AI era and one main advantage of hanamark is how short,compact but full of features. I have carefully built this engine such that even llms can easily understand the whole context of hanamark with limited token and information so that they can assist you to setup and your sites. This page contains the complete Hanamark documentation in a format optimized for Large Language Models (ChatGPT, Claude, Gemini, etc.). Copy the sections below to give an LLM context about Hanamark.

## Quick Context Prompt

Copy this to quickly give an LLM context about Hanamark:

```
Hanamark is a static site generator written in Go. Here's what you need to know:

COMMANDS:
- hanamark init    # Initialize new project
- hanamark build   # Build static site
- hanamark serve   # Start dev server on port 3000

PROJECT STRUCTURE:
configurables/
  config.json         # Main configuration
  source_md/          # Markdown content
  templates/          # HTML templates (Go html/template)
  static/             # CSS, JS, images -> copied to /static/

KEY CONCEPTS:
- single.html: Template for individual pages (.md files)
- list.html: Template for section index pages (folders with _index.md)
- _base.html: Master layout wrapping all pages
- _header.html, _footer.html: Partials included in _base.html

FRONT MATTER (YAML at top of .md files):
---
created_on: 2024-01-15    # Required for single pages
tags: ["tag1", "tag2"]    # Optional
draft: true               # Exclude from build
template: "custom.html"   # Override template
rss: true                 # for _index.md only ie on list pages
---

TEMPLATE VARIABLES:
- {{ .GenHtml }}      # Rendered markdown content
- {{ .PageTitle }}    # First heading from content
- {{ .CreatedDate }}  # Creation date (use .Format "Jan 2, 2006")
- {{ .Tags }}         # Array of tags
- {{ .List }}         # Array of pages (in list.html)

CONFIG.JSON STRUCTURE:
{
  "filepath": {
    "sourceMDDir": "./configurables/source_md",
    "destHtmlDir": "./output_html",
    "templatePath": "./configurables/templates",
    "sourceStaticFiles": "./configurables/static"
  },
  "tags": true,
  "servePort": "3000"
}
```

## Command Reference

All available CLI commands:

```bash
# Initialize a new Hanamark project
./hanamark init

# Build the static site (generates output_html/)
./hanamark build

# Start local development server
./hanamark serve

# Show help
./hanamark help

# Show version
./hanamark -version

# Typical workflow
./hanamark init
# Edit content in configurables/source_md/
./hanamark build
./hanamark serve
# Deploy output_html/ to hosting
```

## Template Examples

Common template patterns:

```html
<!-- _base.html - Master layout -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <link rel="stylesheet" href="/static/css/styles.css">
  <title>{{ .PageTitle }}</title>
</head>
<body>
  {{ template "_header.html" }}
  {{ block "main" . }}{{ end }}
  {{ template "_footer.html" }}
</body>
</html>

<!-- single.html - Individual page template -->
<article>
  <h1>{{ .PageTitle }}</h1>
  <time>{{ .CreatedDate.Format "January 2, 2006" }}</time>
  {{ range .Tags }}
    <a href="{{ .TagDestPath }}">{{ .TagName }}</a>
  {{ end }}
  <div>{{ .GenHtml }}</div>
</article>

<!-- list.html - Section listing template -->
<section>
  <h1>{{ .PageTitle }}</h1>
  <ul>
    {{ range .List }}
      <li>
        <a href="{{ .DestPageDir }}">{{ .PageTitle }}</a>
        <time>{{ .CreatedDate.Format "Jan 2, 2006" }}</time>
      </li>
    {{ end }}
  </ul>
</section>
```

## Configuration Reference

Complete config.json with all options:

```json
{
  "filepath": {
    "sourceMDDir": "./configurables/source_md",
    "destHtmlDir": "./output_html",
    "templatePath": "./configurables/templates",
    "sourceStaticFiles": "./configurables/static",
    "mdAssetsSourcePath": "./configurables/source_md/assets/",
    "mdAssetsDestPath": "./output_html/assets/"
  },
  "indexHomepageHtml": {
    "type": "section",
    "name": "blog"
  },
  "logger": {
    "filepath": "hanamark.logs",
    "level": "debug"
  },
  "rss": {
    "isRssEnabled": true,
    "title": "My Blog",
    "link": "https://example.com",
    "authorName": "Your Name",
    "authorEmailID": "you@example.com",
    "rssOutputName": "feed.xml"
  },
  "tags": true,
  "sortFilesByCreatedOn": true,
  "servePort": "3000"
}
```

## Full Documentation

Complete README with all details - see the full [Hanamark README](https://github.com/thevoid12/hanamark) for comprehensive documentation including:

- Features
- Installation (Download binary or build from source)
- Quick Start Guide
- Project Structure
- Mandatory Requirements
- Configuration Options
- Content Organization
- Front Matter Fields
- Templating System
- Tags System
- Assets and Static Files
- RSS Feed Generation
- CLI Commands
- Examples

For the complete documentation, visit the [Hanamark GitHub repository](https://github.com/thevoid12/hanamark).
