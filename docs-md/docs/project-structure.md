---
created_on: 2024-12-26
---

# Project Structure

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
│       ├── _opengraph.html      # OpenGraph meta tags partial
│       ├── single.html          # Default single page template
│       ├── blog/
│       │   ├── single.html      # Blog post template
│       │   └── list.html        # Blog listing template
│       └── tags/
│           ├── single.html      # All tags page
│           └── list.html        # Individual tag page
└── output_html/                 # Generated site (after build)
```

## Key Directories

### `source_md/`

Contains all your Markdown content. The directory structure here is mirrored in the output.

### `templates/`

HTML templates using Go's `html/template` syntax. Files starting with `_` are system files (base, header, footer).

### `static/`

Static assets copied directly to `output_html/static/`. Access them via `/static/` in your templates.

### `output_html/`

Generated after running `hanamark build`. This is your deployable static site.

---

## Mandatory Requirements

These are **required** for Hanamark to work correctly:

### Required Directory Structure

```
configurables/
├── config.json                    # REQUIRED
├── source_md/                     # REQUIRED
│   └── assets/                    # REQUIRED - All markdown images here
├── static/                        # REQUIRED
│   └── css/
│       └── styles.css             # REQUIRED - This exact path
└── templates/                     # REQUIRED
    ├── _base.html                 # REQUIRED
    ├── _header.html               # REQUIRED (can be empty)
    ├── _footer.html               # REQUIRED (can be empty)
    └── single.html                # REQUIRED
```

### Required Files

| File | Path | Can Be Empty? |
|------|------|---------------|
| `config.json` | `configurables/config.json` | No |
| `_base.html` | `templates/_base.html` | No |
| `_header.html` | `templates/_header.html` | Yes |
| `_footer.html` | `templates/_footer.html` | Yes |
| `single.html` | `templates/single.html` | No |
| `styles.css` | `static/css/styles.css` | Yes |

### Required Paths (Cannot Be Changed)

| Asset Type | Must Be At | Referenced As |
|------------|------------|---------------|
| CSS | `static/css/styles.css` | `/static/css/styles.css` |
| Favicon | `static/favicon.ico` | `/static/favicon.ico` |
| Images | `static/images/` | `/static/images/` |
| Markdown Assets | `source_md/assets/` | `./assets/` in markdown |

---

## System Files (Underscore Prefix)

Files starting with `_` are **system files** with special meaning:

| File | Location | Purpose | Required? |
|------|----------|---------|-----------|
| `_base.html` | `templates/` | Master layout template | **Yes** |
| `_header.html` | `templates/` | Header partial | **Yes** (can be empty) |
| `_footer.html` | `templates/` | Footer partial | **Yes** (can be empty) |
| `_opengraph.html` | `templates/` | OpenGraph meta tags | No (optional) |
| `_index.html` | `templates/` | Custom homepage template | No |
| `_index.md` | Any content folder | Section config for list pages | No (but required for lists) |

> **Warning:** Do not delete `_header.html` or `_footer.html` even if empty. The `_base.html` template references them and will fail if they don't exist.
