---
created_on: 2024-12-26
---

# Front Matter

Front matter is YAML metadata at the top of Markdown files, enclosed by `---`.

## Basic Example

```markdown
---
created_on: 2024-01-15
tags: ["go", "tutorial"]
draft: false
template: "custom_template.html"
rss: true
title: "example"
description: "A comprehensive guide to OpenGraph in Hanamark"
ogImage: "/static/images/opengraph-guide.png"
---

# Your Content Here
```

## Supported Fields

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `created_on` | string | Creation date | **Yes** (for single pages) |
| `tags` | array | List of tags | No |
| `draft` | boolean | If `true`, page is excluded from build | No |
| `template` | string | Custom template path (relative to templates dir) | No |
| `rss` | boolean | Include section in RSS feed (for `_index.md` only) | No |
| `title` | string | Custom title for document (overrides H1) | No |
| `description` | string | Description for meta tags (OpenGraph, Twitter) | No |
| `ogImage` | string | Custom image URL for social media previews | No |

## Date Formats

Hanamark supports multiple date formats:

### ISO 8601

```yaml
created_on: 2024-01-15
created_on: 2024-01-15T10:30:00Z
created_on: 2024-01-15 10:30:45
```

### Slash-separated

```yaml
created_on: 2024/01/15
```

### Human-friendly

```yaml
created_on: 15-01-2024
created_on: 15/01/2024
created_on: 15 Jan 2024
created_on: 15 January 2024
```

### With Time

```yaml
created_on: 15 Jan 2024 10:30
created_on: 15 January 2024 10:30
```

### Compact

```yaml
created_on: 20240115
```

## Draft Mode

Mark a post as draft to exclude it from the build:

```markdown
---
created_on: 2024-01-20
draft: true
---

# Work in Progress

This post won't appear in the build.
```

## Custom Templates

Override the default template for a specific page:

```markdown
---
created_on: 2024-01-25
template: "blog/featured_post.html"
---

# Featured Article

This uses a custom template.
```
