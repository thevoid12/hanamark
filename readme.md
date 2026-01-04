# Hanamark

**Hanamark** is a fast, minimal, and flexible static site generator written in Go. I wrote this engine for my personal blog [thisisvoid.in](https://www.thisisvoid.in) and now putting it out so that anyone can quickly build and maintain their static sites.

> Write in Markdown. Style with templates. Deploy anywhere.
---
- contributions are highly welcomed. If you find any issue please report the issues I will resolve it quickly.
- for a quick setup video demo: https://drive.google.com/file/d/1tQX_J_ieNC1rQxcFhsn7sqlQlyvOfdbE/view?usp=sharing
- llm friendly doc (click the copy button,paste it in your fav llm to ease setup): https://hanamark.thisisvoid.in/llm-docs 
---

## Table of Contents

- [Hanamark](#hanamark)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Quick Start](#quick-start)
  - [Installation](#installation)
    - [Download Binary (Recommended)](#download-binary-recommended)
    - [From Source (Optional)](#from-source-optional)
  - [Project Structure](#project-structure)
  - [Mandatory Requirements](#mandatory-requirements)
    - [Required Directory Structure](#required-directory-structure)
    - [Required Files](#required-files)
    - [Required Paths (Cannot Be Changed)](#required-paths-cannot-be-changed)
    - [System Files (Underscore Prefix)](#system-files-underscore-prefix)
  - [Rules and Constraints](#rules-and-constraints)
    - [What You CAN Do](#what-you-can-do)
    - [What You CANNOT Do](#what-you-cannot-do)
    - [Template Rules](#template-rules)
    - [Content Rules](#content-rules)
  - [Configuration](#configuration)
    - [Complete Configuration Example](#complete-configuration-example)
    - [File Paths](#file-paths)
    - [Index Homepage](#index-homepage)
    - [RSS Feed](#rss-feed)
    - [Logger](#logger)
    - [Other Options](#other-options)
  - [Content Organization](#content-organization)
    - [Directory Structure](#directory-structure)
    - [Single Pages](#single-pages)
    - [List Pages (Sections)](#list-pages-sections)
  - [Front Matter](#front-matter)
    - [Supported Fields](#supported-fields)
    - [Date Formats](#date-formats)
  - [Templating](#templating)
    - [Understanding single.html vs list.html](#understanding-singlehtml-vs-listhtml)
    - [System Files (Underscore-Prefixed Files)](#system-files-underscore-prefixed-files)
    - [Template Types](#template-types)
    - [The \_base.html Template (Global Layout)](#the-_basehtml-template-global-layout)
    - [Content Templates](#content-templates)
    - [Template Variables](#template-variables)
      - [Single Page Variables (`.PageMeta`)](#single-page-variables-pagemeta)
      - [List Page Variables](#list-page-variables)
      - [Tag List Variables](#tag-list-variables)
    - [Custom Templates](#custom-templates)
    - [Template Lookup Order](#template-lookup-order)
  - [Tags](#tags)
    - [Tag Templates](#tag-templates)
  - [Assets and Static Files](#assets-and-static-files)
    - [Static Files](#static-files)
    - [Markdown Assets](#markdown-assets)
  - [RSS Feed Generation](#rss-feed-generation)
    - [Setup Steps](#setup-steps)
  - [Commands](#commands)
    - [`hanamark init`](#hanamark-init)
    - [`hanamark build`](#hanamark-build)
    - [`hanamark serve`](#hanamark-serve)
    - [`hanamark help`](#hanamark-help)
    - [`hanamark -version`](#hanamark--version)
  - [Examples](#examples)
    - [Basic Blog Post](#basic-blog-post)
    - [Draft Post](#draft-post)
    - [Custom Template Post](#custom-template-post)
    - [Section Configuration](#section-configuration)
    - [Complete Template Example](#complete-template-example)
    - [Adding JavaScript to Templates](#adding-javascript-to-templates)
  - [Building from Source](#building-from-source)
    - [Prerequisites](#prerequisites)
    - [Build](#build)
    - [Development](#development)
  - [License](#license)
  - [Contributing](#contributing)

---

## Features

- **Markdown to HTML** - Write content in Markdown, get beautiful HTML
- **Front Matter Support** - YAML front matter for metadata (tags, dates, custom templates, draft mode)
- **Template Inheritance** - Base templates with partials (`_base.html`, `_header.html`, `_footer.html`)
- **Automatic List Pages** - Sections automatically generate index pages with links to content
- **Tag System** - Organize content with tags, auto-generated tag pages
- **RSS Feed** - Built-in RSS  feed generation
- **Draft Mode** - Mark posts as drafts to exclude from build
- **Flexible Date Parsing** - Supports multiple date formats
- **Static Assets** - CSS, images, and other static files copied automatically
- **Local Development Server** - Built-in server for previewing your site
- **Zero JavaScript by Default** - Pure static HTML output, but you can add JavaScript via templates

---

## Quick Start

```bash
# 1. Download the binary or build hanamark binary
go build -o hanamark .

# 2. Initialize a new project
./hanamark init

# 3. Build the site
./hanamark build

# 4. Preview locally
./hanamark serve
```

By default (can be configured in config.json) if served,

 Your site will be available at

`http://localhost:3000` default or whatever port configured in config.json

---

## Installation

### Download Binary (Recommended)

Download the pre-built binary for your platform from the [GitHub Releases](https://github.com/thevoid12/hanamark/releases) page.

| Platform | Download |
| --- | --- |
| macOS (Intel) | `hanamark_Darwin_x86_64.tar.gz` |
| macOS (Apple Silicon) | `hanamark_Darwin_arm64.tar.gz` |
| Linux (64-bit) | `hanamark_Linux_x86_64.tar.gz` |
| Linux (ARM 64-bit) | `hanamark_Linux_arm64.tar.gz` |
| Windows (64-bit) | `hanamark_Windows_x86_64.zip` |
| Windows (ARM 64-bit) | `hanamark_Windows_arm64.zip` |

download the binary and use gunzip to extract the tarball and then get the executable

```bash
# Example: Download and make executable (macOS/Linux)
chmod +x hanamark-darwin-arm64
mv hanamark-darwin-arm64 hanamark

# Run
./hanamark init
```

> **No dependencies required** - the binary is self-contained.

### From Source (Optional)

Only needed if you want to build from source or contribute to development.

**Prerequisites:** Go 1.20 or higher

```bash
git clone https://github.com/thevoid12/hanamark.git
cd hanamark
go build -o hanamark .
```

Or using Make:

```bash
make build
```

---

## Project Structure

After running `hanamark init`, you get the following structure:

```
your-project/
├── configurables/
│   ├── config.json              # Main configuration file
│   ├── source_md/               # Your Markdown content
│   │   ├── index.md             # Homepage (optional)
│   │   ├── about.md             # Single page example
│   │   ├── blog/                # Section with multiple posts
│   │   │   ├── _index.md        # Section configuration
│   │   │   ├── post-1.md
│   │   │   └── post-2.md
│   │   └── assets/              # Images used in markdown
│   │       └── image.png
│   ├── static/                  # Static files (CSS, favicon, etc.)
│   │   ├── css/
│   │   │   └── styles.css
│   │   ├── images/
│   │   └── favicon.ico
│   └── templates/               # HTML templates
│       ├── _base.html           # Base layout template
│       ├── _header.html         # Header partial
│       ├── _footer.html         # Footer partial
│       ├── single.html          # Default single page template
│       ├── blog/
│       │   ├── single.html      # Blog post template
│       │   └── list.html        # Blog listing template
│       └── tags/
│           ├── single.html      # All tags page
│           └── list.html        # Individual tag page
└── output_html/                 # Generated site (after build)
```

---

## Mandatory Requirements

These are **required** for Hanamark to work correctly. Do not change these conventions.

### Required Directory Structure

```
configurables/
├── config.json                    # REQUIRED - Main configuration file
├── source_md/                     # REQUIRED - Markdown content directory
│   └── assets/                    # REQUIRED - All markdown images must be here
├── static/                        # REQUIRED - Static files directory
│   └── css/
│       └── styles.css             # REQUIRED - Main stylesheet (this exact path)
└── templates/                     # REQUIRED - HTML templates directory
    ├── _base.html                 # REQUIRED - Base layout template
    ├── _header.html               # REQUIRED - Header partial (can be empty)
    ├── _footer.html               # REQUIRED - Footer partial (can be empty)
    └── single.html                # REQUIRED - Default single page template
```

### Required Files

| File | Path | Purpose | Can Be Empty? |
|------|------|---------|---------------|
| `config.json` | `configurables/config.json` | Main configuration | No |
| `_base.html` | `templates/_base.html` | Base layout wrapper | No |
| `_header.html` | `templates/_header.html` | Header partial | Yes |
| `_footer.html` | `templates/_footer.html` | Footer partial | Yes |
| `single.html` | `templates/single.html` | Default page template | No |
| `styles.css` | `static/css/styles.css` | Main stylesheet | Yes |

### Required Paths (Cannot Be Changed)

| Asset Type | Must Be At | Referenced As |
|------------|------------|---------------|
| CSS | `static/css/styles.css` | `/static/css/styles.css` |
| Favicon | `static/favicon.ico` or `static/favicon.svg` | `/static/favicon.ico` |
| Images | `static/images/` | `/static/images/` |
| Markdown Assets | `source_md/assets/` | `./assets/` in markdown |

### System Files (Underscore Prefix)

Files starting with `_` are **system files** with special meaning:

| File | Location | Purpose | Required? |
|------|----------|---------|-----------|
| `_base.html` | `templates/` | Master layout template | **Yes** |
| `_header.html` | `templates/` | Header partial included in base | **Yes** (can be empty) |
| `_footer.html` | `templates/` | Footer partial included in base | **Yes** (can be empty) |
| `_index.html` | `templates/` | Custom homepage template | No |
| `_index.md` | Any content folder | Section configuration for list pages | No (but required for list pages) |

> **Warning:** Do not delete `_header.html` or `_footer.html` even if empty. The `_base.html` template references them and will fail if they don't exist.

---

## Rules and Constraints

### What You CAN Do

- **Add JavaScript** - Include `<script>` tags in `_base.html` or any template
- **Add external CSS** - Link to CDN stylesheets in `_base.html`
- **Create custom templates** - Add any `.html` files in `templates/`
- **Nest content folders** - Create deep folder structures in `source_md/`
- **Use custom front matter** - Add any YAML fields, access via `.FrontMatterMap`
- **Override templates per-page** - Use `template:` in front matter
- **Empty system files** - `_header.html` and `_footer.html` can be empty

### What You CANNOT Do

- **Change static path** - Assets must be in `static/` and referenced as `/static/`
- **Change CSS path** - Must be `static/css/styles.css`
- **Delete system files** - `_base.html`, `_header.html`, `_footer.html` must exist
- **Use RSS on single pages** - RSS only works on sections (list pages with `_index.md`)
- **Skip `_index.md` for lists** - Folders without `_index.md` won't generate list pages
- **Use relative paths for static assets** - Always use absolute paths like `/static/...`

### Template Rules

1. **`_base.html` must contain:**
   ```html
   {{ template "_header.html" }}
   {{ block "main" . }}{{ end }}
   {{ template "_footer.html" }}
   ```

2. **Content templates are auto-wrapped** - Don't add `{{ define "main" }}` manually

3. **Template lookup order:**
   - Section-specific: `templates/blog/single.html`
   - Root fallback: `templates/single.html`

### Content Rules

1. **Single pages require `created_on`** - Front matter must include a date
2. **List pages require `_index.md`** - Even if empty, the file must exist
3. **Assets must be in one folder** - All markdown images go in `source_md/assets/`
4. **Reference assets relatively** - Use `./assets/image.png` in markdown

---

## Configuration

Configuration is stored in `./configurables/config.json`.

### Complete Configuration Example

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

### File Paths

| Key | Description | Example |
|-----|-------------|---------|
| `sourceMDDir` | Directory containing Markdown source files | `./configurables/source_md` |
| `destHtmlDir` | Output directory for generated HTML | `./output_html` |
| `templatePath` | Directory containing HTML templates | `./configurables/templates` |
| `sourceStaticFiles` | Static files (CSS, JS, images) to copy | `./configurables/static` |
| `mdAssetsSourcePath` | Assets referenced in Markdown files | `./configurables/source_md/assets/` |
| `mdAssetsDestPath` | Destination for Markdown assets | `./output_html/assets/` |

### Index Homepage

Controls what content appears on the root `index.html`:

```json
"indexHomepageHtml": {
  "type": "section",
  "name": "blog"
}
```

| Type | Behavior |
|------|----------|
| `section` | Uses the specified section's list as the homepage |
| `page` | Copies the specified page as `index.html` |

**Examples:**

```json
// Use blog section as homepage
// Uses list.html -> generates blog/index.html with list of posts
"indexHomepageHtml": {
  "type": "section",
  "name": "blog"
}

// Use about.html as homepage
"indexHomepageHtml": {
  "type": "page",
  "name": "about.html"
}
```

**Nested Folder Example:**

```
source_md/
└── blog/
    ├── _index.md
    ├── post-1.md
    └── tutorials/           # Nested folder
        ├── _index.md
        └── go-basics.md
```

To use the nested `tutorials` section as homepage:

```json
// For nested folders, use path from source_md root
"indexHomepageHtml": {
  "type": "section",
  "name": "blog/tutorials"
}

// Or if tutorials is at root level
"indexHomepageHtml": {
  "type": "section",
  "name": "tutorials"
}
```

**How List Pages Work:**

```
blog/
├── _index.md      -> Uses list.html   -> blog/index.html (shows links to all posts)
├── post-1.md      -> Uses single.html -> blog/post-1.html
└── post-2.md      -> Uses single.html -> blog/post-2.html
```

> **Important:** If there is no `_index.md` in a folder, no list page (`index.html`) will be generated. The engine will treat files as individual pages only.

**Custom Homepage Template (`_index.html`):**

When using `type: "section"`, you can create a custom template for the root `index.html` by adding `_index.html` in your templates root directory:

```
templates/
├── _index.html    # Custom homepage template (optional)
├── _base.html
├── single.html
└── blog/
    └── list.html
```

- If `_index.html` exists in templates root, it will be used for the homepage
- If `_index.html` does not exist, the referenced section's template (e.g., `blog/list.html`) will be used

> **Note:** If you have an `index.md` in your source root, it takes precedence over all these settings.

### RSS Feed

```json
"rss": {
  "isRssEnabled": true,
  "title": "My Blog",
  "link": "https://example.com",
  "authorName": "Your Name",
  "authorEmailID": "you@example.com",
  "rssOutputName": "feed.xml"
}
```

| Key | Description | Required |
|-----|-------------|----------|
| `isRssEnabled` | Enable/disable RSS generation | Yes |
| `title` | Feed title | Yes |
| `link` | Your site's base URL | Yes |
| `authorName` | Author name for feed items | No |
| `authorEmailID` | Author email | No |
| `rssOutputName` | Output filename (default: `feed.xml`) | No |

### OpenGraph Support (Social Media Metadata)

Hanamark supports OpenGraph and Twitter Card metadata generation. This allows your pages to look great when shared on social media.

**Configuration:**

```json
"opengraph": {
  "enabled": true,
  "siteName": "My Awesome Blog",
  "baseUrl": "https://example.com",
  "defaultImage": "/static/images/og-default.png",
  "imageWidth": "1200",
  "imageHeight": "630",
  "twitterCard": "summary_large_image"
}
```

| Key | Description | Required |
|-----|-------------|----------|
| `enabled` | Enable/disable generation | Yes |
| `siteName` | Site name for `og:site_name` | No |
| `baseUrl` | Base URL (essential for valid absolute URLs) | **Yes** |
| `defaultImage` | Default fallback image if page has no image | No |
| `twitterCard` | Twitter card type (usually `summary_large_image`) | No |

**Front Matter Overrides:**

You can override OpenGraph data per page using standard front matter fields:

```markdown
---
title: "My Special Post"
description: "Custom description for social previews"
ogImage: "/static/images/special-cover.png"
---
```

- `description`: Used for `og:description` and `twitter:description`
- `ogImage`: Used for `og:image` and `twitter:image` (overrides `defaultImage`)

### Logger

```json
"logger": {
  "filepath": "hanamark.logs",
  "level": "debug"
}
```

| Level | Description |
|-------|-------------|
| `debug` | Verbose logging for development |
| `info` | General information |
| `warn` | Warnings only |
| `error` | Errors only |

### Other Options

| Key | Type | Description |
|-----|------|-------------|
| `tags` | boolean | Enable tag system |
| `sortFilesByCreatedOn` | boolean | Sort lists by `created_on` date (true) or `updated` date (false) |
| `servePort` | string | Port for local development server (default: `3000`) |

---

## Content Organization

### Directory Structure

Hanamark mirrors your source directory structure in the output:

```
source_md/                    output_html/
├── about.md            ->    ├── about.html
├── projects.md         ->    ├── projects.html
└── blog/                     └── blog/
    ├── _index.md       ->        ├── index.html (list page)
    ├── post-1.md       ->        ├── post-1.html
    └── post-2.md       ->        └── post-2.html
```

### Single Pages

Any `.md` file (except `_index.md`) becomes a single HTML page:

```markdown
---
created_on: 2024-01-15
tags: ["about", "personal"]
---

# About Me

This is my about page content...
```

**Output:** `about.html`

### List Pages (Sections)

Directories with a `_index.md` file become **sections** with automatic list pages if you want list page it is mandatory to have _index.md in that folder:

```
blog/
├── _index.md           # Section configuration
├── post-1.md           # Blog post
├── post-2.md           # Blog post
└── nested/             # Nested section
    ├── _index.md
    └── deep-post.md
```

The `_index.md` file configures the section:

```markdown
---
rss: true
template: "blog/custom_list_template.html"
---
<!-- Optional content or comments -->
```

**Output:** `blog/index.html` with links to all posts in the section.

---

## Front Matter

Front matter is YAML metadata at the top of Markdown files, enclosed by `---`:

```markdown
---
created_on: 2024-01-15
tags: ["go", "tutorial"]
draft: false
template: "custom_template.html"
rss: true
---

# Your Content Here
```

### Supported Fields

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `created_on` | string | Creation date | **Yes** (for single pages) |
| `tags` | array | List of tags | No |
| `draft` | boolean | If `true`, page is excluded from build | No |
| `template` | string | Custom template path (relative to templates dir) | No |
| `rss` | boolean | Include section in RSS feed (for `_index.md` only) | No |

### Date Formats

Hanamark supports multiple date formats:

```yaml
# ISO 8601
created_on: 2024-01-15
created_on: 2024-01-15T10:30:00Z
created_on: 2024-01-15 10:30:45

# Slash-separated
created_on: 2024/01/15

# Human-friendly
created_on: 15-01-2024
created_on: 15/01/2024
created_on: 15 Jan 2024
created_on: 15 January 2024

# With time
created_on: 15 Jan 2024 10:30
created_on: 15 January 2024 10:30

# Compact
created_on: 20240115
```

---

## Templating

Hanamark uses Go's `html/template` package with a hierarchical template system.

### Understanding single.html vs list.html

This is the most important concept in Hanamark templating:

| Template | What It Does | When It's Used |
|----------|--------------|----------------|
| **`single.html`** | Renders **individual content pages** | Used for every `.md` file (e.g., `about.md` -> `about.html`, `blog/post-1.md` -> `blog/post-1.html`) |
| **`list.html`** | Renders **index/listing pages** with links to content | Used for sections with `_index.md` (e.g., `blog/_index.md` -> `blog/index.html` showing list of all posts) |

**Think of it this way:**
- `single.html` = Template for **reading a single article/page**
- `list.html` = Template for **browsing a collection of articles** (like a blog archive)

**Example:**
```
blog/
├── _index.md      -> Uses list.html   -> blog/index.html (shows links to all posts)
├── post-1.md      -> Uses single.html -> blog/post-1.html (the actual post content)
└── post-2.md      -> Uses single.html -> blog/post-2.html (the actual post content)
```

### System Files (Underscore-Prefixed Files)

Files starting with `_` are **system files** used by Hanamark. **Do not delete them.**

| File | Purpose | Can You Empty It? |
|------|---------|-------------------|
| `_base.html` | Base layout wrapper for all pages | No - required for proper HTML structure |
| `_header.html` | Site header/navigation included in every page | Yes - if you don't want a header |
| `_footer.html` | Site footer included in every page | Yes - if you don't want a footer |
| `_index.html` | Custom homepage template (in templates root) | Optional - only create if you want a custom homepage layout |
| `_index.md` | Section configuration (in content folders) | Yes - but keep the file to enable list pages |

> **Important:** You can empty `_header.html` or `_footer.html` if you don't need them, but keep the files. Deleting them will cause build errors since `_base.html` references them.

### Template Types

| Template | Purpose | Location |
|----------|---------|----------|
| `_base.html` | Base layout with `<html>`, `<head>`, `<body>` | Root of templates |
| `_header.html` | Site header/navigation | Root of templates |
| `_footer.html` | Site footer | Root of templates |
| `single.html` | Individual content pages | Root or section folder |
| `list.html` | Section listing pages | Section folder |

### The _base.html Template (Global Layout)

`_base.html` is the **master template** that wraps every page on your site. Use it for:

- **HTML document structure** (`<!DOCTYPE>`, `<html>`, `<head>`, `<body>`)
- **Global CSS/JS** - stylesheets and scripts needed on every page
- **Meta tags** - charset, viewport, SEO tags
- **Favicon** - site icon
- **Common layout** - header, footer, navigation

Everything in `_base.html` appears on **every single page** of your site.

**`_base.html`** (Base Layout):
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="icon" href="/static/favicon.ico">
  <link rel="stylesheet" href="/static/css/styles.css">
  <title>{{ .PageTitle }}</title>
</head>
<body>
  {{ template "_header.html" }}

  {{ block "main" . }}{{ end }}

  {{ template "_footer.html" }}
</body>
</html>
```

**Key points:**
- `{{ template "_header.html" }}` - Includes the header partial
- `{{ block "main" . }}{{ end }}` - This is where `single.html` or `list.html` content gets injected
- `{{ template "_footer.html" }}` - Includes the footer partial
- Static assets are always served from `/static/` (e.g., `/static/css/styles.css`, `/static/favicon.ico`, `/static/images/logo.png`)

**`_header.html`** (Header Partial):
```html
<header>
  <nav>
    <a href="/">Home</a>
    <a href="/about.html">About</a>
    <a href="/blog/index.html">Blog</a>
    <a href="/feed.xml">RSS</a>
  </nav>
</header>
```

**`_footer.html`** (Footer Partial):
```html
<footer>
  <p>Built with Hanamark</p>
</footer>
```

### Content Templates

**`single.html`** (For individual pages):
```html
<article class="post">
  <span class="date">{{ .CreatedDate.Format "2 January 2006" }}</span>
  <div class="content">
    {{ .GenHtml }}
  </div>
</article>
```

**`list.html`** (For section index pages):
```html
<section>
  <h1>{{ .PageTitle }}</h1>
  <ul>
    {{ range .List }}
      <li><a href="{{ .DestPageDir }}">{{ .PageTitle }}</a></li>
    {{ end }}
  </ul>
</section>
```

> **Note:** Content templates are automatically wrapped with `{{ define "main" }}...{{ end }}` if `_base.html` exists.

### Template Variables

#### Single Page Variables (`.PageMeta`)

| Variable | Type | Description |
|----------|------|-------------|
| `.GenHtml` | string | Rendered HTML content from Markdown |
| `.PageTitle` | string | First heading extracted from content |
| `.PageName` | string | Page name |
| `.CreatedDate` | time.Time | Creation date from front matter |
| `.UpdatedDate` | time.Time | Last modification time of file |
| `.DestPageDir` | string | Destination path |
| `.ReadTime` | int | Estimated reading time |
| `.FrontMatterMap` | map | All front matter as key-value pairs |
| `.Tags` | []*Tag | List of tags |

**Example Usage:**
```html
<article>
  <h1>{{ .PageTitle }}</h1>
  <time>{{ .CreatedDate.Format "January 2, 2006" }}</time>
  <time>Updated: {{ .UpdatedDate.Format "2006-01-02" }}</time>

  <div class="tags">
    {{ range .Tags }}
      <a href="{{ .TagDestPath }}">{{ .TagName }}</a>
    {{ end }}
  </div>

  <div class="content">
    {{ .GenHtml }}
  </div>
</article>
```

#### List Page Variables

| Variable | Type | Description |
|----------|------|-------------|
| `.PageTitle` | string | Section/folder name |
| `.List` | []*PageMeta | Array of pages in section |

**Example Usage:**
```html
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

#### Tag List Variables

For `tags/single.html` (main tags page):

| Variable | Type | Description |
|----------|------|-------------|
| `.List` | []*TagList | All tags with counts |
| `.List[].TagName` | string | Tag name |
| `.List[].TagDestPath` | string | Path to tag page |
| `.List[].Count` | int | Number of posts with tag |

```html
<h1>All Tags</h1>
<ul>
  {{ range .List }}
    <li>
      <a href="{{ .TagDestPath }}">{{ .TagName }}</a>
      ({{ .Count }} posts)
    </li>
  {{ end }}
</ul>
```

For `tags/list.html` (individual tag page):

| Variable | Type | Description |
|----------|------|-------------|
| `.PageTitle` | string | Tag name |
| `.List` | []*Tag | Posts with this tag |
| `.List[].FileHeading` | string | Post title |
| `.List[].FileDestPath` | string | Path to post |

```html
<h1>Posts tagged "{{ .PageTitle }}"</h1>
<ul>
  {{ range .List }}
    <li><a href="{{ .FileDestPath }}">{{ .FileHeading }}</a></li>
  {{ end }}
</ul>
```

### Custom Templates

Override templates per-page using front matter:

```markdown
---
created_on: 2024-01-15
template: "blog/custom_post.html"
---

# Special Post

This uses a custom template!
```

For sections, specify in `_index.md`:

```markdown
---
template: "blog/custom_list_template.html"
---
```

### Template Lookup Order

For single pages, Hanamark searches upward from the content's location:

1. `templates/blog/single.html` (section-specific)
2. `templates/single.html` (root fallback)

For list pages:

1. Custom template from `_index.md` front matter
2. `templates/blog/list.html` (section-specific)
3. `templates/list.html` (root fallback)

---

## Tags

Enable tags in config:

```json
{
  "tags": true
}
```

Add tags to content:

```markdown
---
created_on: 2024-01-15
tags: ["go", "tutorial", "beginner"]
---
```

Hanamark generates:

- `output_html/tags/index.html
` - List of all tags
- `output_html/tags/go.html` - Posts tagged "go"
- `output_html/tags/tutorial.html` - Posts tagged "tutorial"

### Tag Templates

Create custom tag templates:

```
templates/
└── tags/
    ├── single.html    # Main tags listing page
    └── list.html      # Individual tag page
```

---

## Assets and Static Files

### Static Files

Files in `sourceStaticFiles` are copied to `destHtmlDir/static/`. **All static assets are always served from the `/static/` path.**

```
configurables/static/          output_html/static/
├── css/                  ->   ├── css/
│   └── styles.css             │   └── styles.css
├── images/               ->   ├── images/
│   └── logo.png               │   └── logo.png
└── favicon.ico           ->   └── favicon.ico
```

**Standard paths in your templates:**

| Asset Type | Path |
|------------|------|
| CSS | `/static/css/styles.css` |
| Favicon | `/static/favicon.ico` |
| Images | `/static/images/your-image.png` |
| JavaScript | `/static/js/script.js` |
| Fonts | `/static/fonts/font.woff2` |

Reference in templates:

```html
<link rel="icon" href="/static/favicon.ico">
<link rel="stylesheet" href="/static/css/styles.css">
<img src="/static/images/logo.png" alt="Logo">
<script src="/static/js/script.js"></script>
```

> **Note:** Always use absolute paths starting with `/static/` to ensure assets load correctly on all pages regardless of URL depth.

### Markdown Assets

Images referenced in Markdown are copied from `mdAssetsSourcePath` to `mdAssetsDestPath`:

```markdown
![My Image](./assets/photo.png)
```

---

## RSS Feed Generation

> **Important:** RSS feeds can **only** be enabled on **list pages** (sections). You must add `rss: true` in the `_index.md` file of the section you want to include in the feed. RSS cannot be enabled on individual single pages.

**How it works:**
1. Enable RSS globally in `config.json`
2. Add `rss: true` to the `_index.md` of sections you want in the feed
3. All posts within those sections are automatically added to the RSS feed

### Setup Steps

1. Enable RSS in config:

```json
"rss": {
  "isRssEnabled": true,
  "title": "My Blog",
  "link": "https://example.com",
  "authorName": "Your Name",
  "authorEmailID": "you@example.com",
  "rssOutputName": "feed.xml"
}
```

2. Mark sections for RSS in `_index.md` (this is required - RSS only works on list pages):

```markdown
---
rss: true
---
```

3. Build your site - `feed.xml` is generated in the output root.

Link to feed in templates:

```html
<link rel="alternate" type="application/rss+xml"
      title="RSS Feed" href="/feed.xml">
```

---

## Commands

### `hanamark init`

Initialize a new project in the current directory:

```bash
./hanamark init
```

Creates a `configurables/` directory with:
- Default `config.json`
- Sample templates
- Example content
- Static file structure

### `hanamark build`

Build the static site:

```bash
./hanamark build
```

- Parses all Markdown files
- Applies templates
- Copies static assets
- Generates tag pages
- Creates RSS feed (if enabled)
- Outputs to `destHtmlDir`

### `hanamark serve`

Start a local development server:

```bash
./hanamark serve
```

Serves the built site at `http://localhost:3000` (or configured port).

### `hanamark help`

Display help information:

```bash
./hanamark help
```

### `hanamark -version`

Show version:

```bash
./hanamark -version
```

---

## Examples

### Basic Blog Post

`source_md/blog/my-first-post.md`:

```markdown
---
created_on: 2024-01-15
tags: ["introduction", "personal"]
---

# My First Post

Welcome to my blog! This is my first post using **Hanamark**.

## What I'll Write About

- Programming tutorials
- Project updates
- Random thoughts

![A nice image](./assets/hero.png)

Thanks for reading!
```

### Draft Post

```markdown
---
created_on: 2024-01-20
draft: true
---

# Work in Progress

This post won't appear in the build until draft is removed or set to false.
```

### Custom Template Post

```markdown
---
created_on: 2024-01-25
template: "blog/featured_post.html"
tags: ["featured"]
---

# Featured Article

This post uses a special template for featured content.
```

### Section Configuration

`source_md/blog/_index.md`:

```markdown
---
rss: true
template: "blog/custom_list_template.html"
---
<!-- This section will be included in RSS and uses a custom list template -->
```

### Complete Template Example

`templates/blog/single.html`:

```html
<article class="blog-post">
  <header>
    <h1>{{ .PageTitle }}</h1>
    <div class="meta">
      <time datetime="{{ .CreatedDate.Format "2006-01-02" }}">
        {{ .CreatedDate.Format "January 2, 2006" }}
      </time>
      {{ if .Tags }}
      <div class="tags">
        {{ range .Tags }}
          <a href="{{ .TagDestPath }}" class="tag">{{ .TagName }}</a>
        {{ end }}
      </div>
      {{ end }}
    </div>
  </header>

  <div class="content">
    {{ .GenHtml }}
  </div>

  <footer>
    <p>Last updated: {{ .UpdatedDate.Format "Jan 2, 2006" }}</p>
  </footer>
</article>
```

### Adding JavaScript to Templates

While Hanamark generates pure static HTML by default, you can easily add JavaScript via templates:

**Example: Adding syntax highlighting with highlight.js**

`templates/_base.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <link rel="stylesheet" href="/static/css/styles.css">
  <!-- External JS library -->
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">
  <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
  <title>{{ .PageTitle }}</title>
</head>
<body>
  {{ template "_header.html" }}
  {{ block "main" . }}{{ end }}
  {{ template "_footer.html" }}

  <!-- Custom JavaScript -->
  <script>
    // Initialize syntax highlighting
    hljs.highlightAll();

    // Add copy buttons to code blocks
    document.querySelectorAll('pre').forEach(function(pre) {
      var btn = document.createElement('button');
      btn.textContent = 'Copy';
      btn.onclick = function() {
        navigator.clipboard.writeText(pre.textContent);
        btn.textContent = 'Copied!';
        setTimeout(function() { btn.textContent = 'Copy'; }, 2000);
      };
      pre.appendChild(btn);
    });
  </script>
</body>
</html>
```

**Example: Dark/Light theme toggle**

```html
<script>
  function toggleTheme() {
    var html = document.documentElement;
    var theme = html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }

  // Load saved theme
  (function() {
    var saved = localStorage.getItem('theme');
    if (saved) document.documentElement.setAttribute('data-theme', saved);
  })();
</script>
```

> **Note:** You can add any JavaScript - analytics, interactive features, external libraries, etc. The "Zero JavaScript by default" means Hanamark itself doesn't require JS, but you're free to add it.

---

## Building from Source

### Prerequisites

- Go 1.20+
- Make (optional)

### Build

```bash
# Clone repository
git clone https://github.com/thevoid12/hanamark.git
cd hanamark

# Build binary
go build -o hanamark .

# Or using Make
make build
```

### Development

```bash
# Run tests
go test ./...

# Build and run
go build -o hanamark . && ./hanamark build
```

---

## License

MIT License - See LICENSE file for details.

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

---

**Built with Hanamark** - A static site generator that gets out of your way.


---
