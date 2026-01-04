---
created_on: 2024-12-26
---

# Tags

Organize your content with tags and auto-generated tag pages.

## Enable Tags

Enable tags in your `config.json`:

```json
{
  "tags": true
}
```

## Add Tags to Content

Add tags to your Markdown files using front matter:

```markdown
---
created_on: 2024-01-15
tags: ["go", "tutorial", "beginner"]
---

# My Tutorial

Content here...
```

## Generated Pages

Hanamark automatically generates:

| Page | Path | Description |
|------|------|-------------|
| All tags | `output_html/tags/tags.html` | List of all tags with counts |
| Tag page | `output_html/tags/go.html` | Posts tagged "go" |
| Tag page | `output_html/tags/tutorial.html` | Posts tagged "tutorial" |

## Tag Templates

Create custom tag templates in `templates/tags/`:

```
templates/
└── tags/
    ├── single.html    # Main tags listing page
    └── list.html      # Individual tag page
```

### Main Tags Page (`tags/single.html`)

```html
<h1>All Tags</h1>
<ul class="tag-list">
  {{ range .List }}
    <li>
      <a href="{{ .TagDestPath }}">{{ .TagName }}</a>
      <span class="count">({{ .Count }})</span>
    </li>
  {{ end }}
</ul>
```

### Individual Tag Page (`tags/list.html`)

```html
<h1>Posts tagged "{{ .PageTitle }}"</h1>
<ul class="post-list">
  {{ range .List }}
    <li>
      <a href="{{ .FileDestPath }}">{{ .FileHeading }}</a>
    </li>
  {{ end }}
</ul>
```

## Linking to Tags

In your templates, link to tag pages:

```html
{{ range .Tags }}
  <a href="{{ .TagDestPath }}" class="tag">{{ .TagName }}</a>
{{ end }}
```

## Best Practices

- Use lowercase tags for consistency
- Use hyphens for multi-word tags: `"web-development"`
- Keep tag names short and descriptive
- Avoid too many tags per post (3-5 is ideal)
