---
created_on: 2024-12-26
---

# Rules and Constraints

This page documents what you can and cannot do with Hanamark. Understanding these constraints will help you avoid common pitfalls.

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

## System Files (Underscore Prefix)

Files starting with `_` are **system files** with special meaning:

| File | Location | Purpose | Required? |
|------|----------|---------|-----------|
| `_base.html` | `templates/` | Master layout template that wraps all pages | **Yes** |
| `_header.html` | `templates/` | Header partial included in base | **Yes** (can be empty) |
| `_footer.html` | `templates/` | Footer partial included in base | **Yes** (can be empty) |
| `_index.html` | `templates/` | Custom homepage template | No |
| `_index.md` | Any content folder | Section configuration for list pages | No (but required for list pages) |

> **Warning:** Do not delete `_header.html` or `_footer.html` even if empty. The `_base.html` template references them and will fail if they don't exist.

---

## What You CAN Do

### Add JavaScript

Hanamark generates pure static HTML by default, but you can add any JavaScript via templates:

```html
<!-- In _base.html -->
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

> **Note:** You can add any JavaScript - analytics, interactive features, external libraries, etc. The "Zero JavaScript by default" means Hanamark itself doesn't require JS to function, but you're completely free to add it via templates.

### Other Things You CAN Do

- **Add external CSS** - Link to CDN stylesheets in `_base.html`
- **Create custom templates** - Add any `.html` files in `templates/`
- **Nest content folders** - Create deep folder structures in `source_md/`
- **Use custom front matter** - Add any YAML fields, access via `.FrontMatterMap`
- **Override templates per-page** - Use `template:` in front matter
- **Empty system files** - `_header.html` and `_footer.html` can be empty

---

## What You CANNOT Do

- **Change static path** - Assets must be in `static/` and referenced as `/static/`
- **Change CSS path** - Must be `static/css/styles.css`
- **Delete system files** - `_base.html`, `_header.html`, `_footer.html` must exist
- **Use RSS on single pages** - RSS only works on sections (list pages with `_index.md`)
- **Skip `_index.md` for lists** - Folders without `_index.md` won't generate list pages
- **Use relative paths for static assets** - Always use absolute paths like `/static/...`

---

## Template Rules

### _base.html Must Contain

```html
{{ template "_header.html" }}
{{ block "main" . }}{{ end }}
{{ template "_footer.html" }}
```

### Content Templates Are Auto-Wrapped

Don't add `{{ define "main" }}` manually - Hanamark does this automatically.

### Template Lookup Order

1. Section-specific: `templates/blog/single.html`
2. Root fallback: `templates/single.html`

---

## Content Rules

1. **Single pages require `created_on`** - Front matter must include a date
2. **List pages require `_index.md`** - Even if empty, the file must exist for list generation
3. **Assets must be in one folder** - All markdown images go in `source_md/assets/`
4. **Reference assets relatively** - Use `./assets/image.png` in markdown

---

## Common Mistakes to Avoid

| Mistake | Why It Fails | Fix |
|---------|--------------|-----|
| Deleting `_header.html` | `_base.html` references it | Keep file, can be empty |
| Using `./static/css/styles.css` | Relative paths break on nested pages | Use `/static/css/styles.css` |
| No `_index.md` in folder | List page won't be generated | Add `_index.md` (can be empty with just `---\n---`) |
| Missing `created_on` | Single pages require a date | Add `created_on: YYYY-MM-DD` to front matter |
| Images outside `assets/` | Won't be copied to output | Put all markdown images in `source_md/assets/` |
