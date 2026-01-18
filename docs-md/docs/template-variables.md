---
created_on: 2024-12-26
---

# Template Variables

Reference for all available template variables in Hanamark.

## Single Page Variables

Available in `single.html` templates:

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

### Accessing FrontMatterMap

Custom frontmatter fields are accessible via `.FrontMatterMap` (keys are lowercased):

```html
{{- if .FrontMatterMap.description -}}
  <meta name="description" content="{{ .FrontMatterMap.description }}">
{{- end -}}

{{- if .FrontMatterMap.ogimage -}}
  <meta property="og:image" content="{{ .FrontMatterMap.ogimage }}">
{{- end -}}
```

### Example Usage

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

## List Page Variables

Available in `list.html` templates:

| Variable | Type | Description |
|----------|------|-------------|
| `.PageTitle` | string | Section/folder name |
| `.List` | []*PageMeta | Array of pages in section |

### Example Usage

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

## Tag List Variables

### Main Tags Page (`tags/single.html`)

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

### Individual Tag Page (`tags/list.html`)

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

## Template Functions

Custom functions available in all templates:

| Function | Description | Example |
|----------|-------------|---------|
| `config` | Get value from config.json | `{{ config "site.title" }}` |
| `hasPrefix` | Check if string starts with prefix | `{{ if hasPrefix $url "http" }}` |
| `urlJoin` | Join URL parts safely | `{{ urlJoin $baseUrl $path }}` |

### Config Function Examples

```html
{{- $siteTitle := config "site.title" -}}
{{- $siteUrl := config "site.url" -}}
{{- $ogEnabled := config "opengraph.enabled" -}}

{{ if $ogEnabled }}
  <meta property="og:site_name" content="{{ $siteTitle }}">
{{ end }}
```

### URL Handling Example

```html
{{- $baseUrl := config "site.url" -}}
{{- $image := .FrontMatterMap.ogimage -}}

<!-- Handle both relative and absolute URLs -->
{{ if hasPrefix $image "http" }}
  <meta property="og:image" content="{{ $image }}">
{{ else }}
  <meta property="og:image" content="{{ urlJoin $baseUrl $image }}">
{{ end }}
```

## Date Formatting

Use Go's time format reference: `Mon Jan 2 15:04:05 MST 2006`

| Format | Output |
|--------|--------|
| `"2006-01-02"` | 2024-01-15 |
| `"January 2, 2006"` | January 15, 2024 |
| `"Jan 2, 2006"` | Jan 15, 2024 |
| `"2 January 2006"` | 15 January 2024 |
| `"02/01/2006"` | 15/01/2024 |
