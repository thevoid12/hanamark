---
created_on: 2024-12-26
---

# Custom Templates

Override templates per-page or per-section using front matter.

## Per-Page Custom Template

Specify a custom template in your Markdown file's front matter:

```markdown
---
created_on: 2024-01-15
template: "blog/custom_post.html"
---

# Special Post

This uses a custom template!
```

The template path is relative to your templates directory.

## Per-Section Custom Template

For sections, specify the template in `_index.md`:

```markdown
---
template: "blog/custom_list_template.html"
---
```

## Creating Custom Templates

### Step 1: Create the Template File

Create a new HTML file in your templates directory:

```
templates/
└── blog/
    └── featured_post.html
```

### Step 2: Write the Template

```html
<article class="featured-post">
  <div class="featured-badge">Featured</div>
  <h1>{{ .PageTitle }}</h1>
  <time>{{ .CreatedDate.Format "January 2, 2006" }}</time>

  {{ if .Tags }}
  <div class="tags">
    {{ range .Tags }}
      <a href="{{ .TagDestPath }}">{{ .TagName }}</a>
    {{ end }}
  </div>
  {{ end }}

  <div class="content">
    {{ .GenHtml }}
  </div>
</article>
```

### Step 3: Reference in Front Matter

```markdown
---
created_on: 2024-01-25
template: "blog/featured_post.html"
tags: ["featured"]
---

# My Featured Article
```

## Template Inheritance

Custom templates still inherit from `_base.html`. Your custom template only replaces the `{{ block "main" }}` content.

## Fallback Behavior

If a custom template is not found, Hanamark falls back to:

1. Section-specific template (e.g., `templates/blog/single.html`)
2. Root template (e.g., `templates/single.html`)
