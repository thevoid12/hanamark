---
created_on: 2024-12-26
---

# Template Basics

Hanamark uses Go's `html/template` package with a hierarchical template system.

## Understanding single.html vs list.html

This is the most important concept in Hanamark templating:

| Template | What It Does | When It's Used |
|----------|--------------|----------------|
| `single.html` | Renders individual content pages | Used for every `.md` file |
| `list.html` | Renders index/listing pages | Used for sections with `_index.md` |

**Think of it this way:**
- `single.html` = Template for reading a single article/page
- `list.html` = Template for browsing a collection of articles

### Example

```
blog/
├── _index.md      -> Uses list.html   -> blog/index.html
├── post-1.md      -> Uses single.html -> blog/post-1.html
└── post-2.md      -> Uses single.html -> blog/post-2.html
```

## System Files

Files starting with `_` are system files used by Hanamark. **Do not delete them.**

| File | Purpose | Can You Empty It? |
|------|---------|-------------------|
| `_base.html` | Base layout wrapper for all pages | No |
| `_header.html` | Site header/navigation | Yes |
| `_footer.html` | Site footer | Yes |
| `_index.html` | Custom homepage template (optional) | Optional |
| `_index.md` | Section configuration (in content folders) | Yes, but keep the file |

## The _base.html Template (Global Layout)

`_base.html` is the **master template** that wraps every page on your site. It is **mandatory** and must exist. Use it for:

- **HTML document structure** (`<!DOCTYPE>`, `<html>`, `<head>`, `<body>`)
- **Global CSS/JS** - stylesheets and scripts needed on every page
- **Meta tags** - charset, viewport, SEO tags
- **Favicon** - site icon
- **Common layout** - header, footer, navigation

Everything in `_base.html` appears on **every single page** of your site.

### Why _base.html is Mandatory

1. **Consistent structure** - Ensures all pages have proper HTML document structure
2. **DRY principle** - Write common elements once, use everywhere
3. **Template inheritance** - Content templates inject into the `{{ block "main" . }}` placeholder
4. **Partials support** - Header and footer are included automatically

### Required Structure

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
- `{{ template "_header.html" }}` - Includes the header partial (file must exist, can be empty)
- `{{ block "main" . }}{{ end }}` - **Required** - Where `single.html` or `list.html` content gets injected
- `{{ template "_footer.html" }}` - Includes the footer partial (file must exist, can be empty)
- Static assets are always served from `/static/` (e.g., `/static/css/styles.css`, `/static/favicon.ico`)

### What You Can Add to _base.html

- **External CSS/JS libraries** (CDN links)
- **Custom JavaScript** (analytics, theme toggles, etc.)
- **Meta tags** (SEO, Open Graph, Twitter cards)
- **Web fonts**
- **Any HTML that should appear on every page**

## Content Templates

**single.html** (For individual pages):

```html
<article class="post">
  <span class="date">{{ .CreatedDate.Format "2 January 2006" }}</span>
  <div class="content">
    {{ .GenHtml }}
  </div>
</article>
```

**list.html** (For section index pages):

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

## Template Lookup Order

For single pages, Hanamark searches upward from the content's location:

1. `templates/blog/single.html` (section-specific)
2. `templates/single.html` (root fallback)

For list pages:

1. Custom template from `_index.md` front matter
2. `templates/blog/list.html` (section-specific)
3. `templates/list.html` (root fallback)

## Adding JavaScript to Templates

While Hanamark generates pure static HTML by default, you can easily add JavaScript via templates:

**Example: Adding syntax highlighting with highlight.js**

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
    hljs.highlightAll();
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
