---
created_on: 2024-12-26
---

# RSS Feed

Hanamark has built-in RSS  feed generation.

> **Important:** RSS feeds can only be enabled on list pages (sections). You must add `rss: true` in the `_index.md` file of the section you want to include in the feed.

## How It Works

1. Enable RSS globally in `config.json`
2. Add `rss: true` to the `_index.md` of sections you want in the feed
3. All posts within those sections are automatically added to the RSS feed

## Configuration

Enable RSS in your `config.json`:

```json
{
  "rss": {
    "isRssEnabled": true,
    "title": "My Blog",
    "link": "https://example.com",
    "authorName": "Your Name",
    "authorEmailID": "you@example.com",
    "rssOutputName": "feed.xml"
  }
}
```

### RSS Options

| Key | Description | Required |
|-----|-------------|----------|
| `isRssEnabled` | Enable/disable RSS generation | Yes |
| `title` | Feed title | Yes |
| `link` | Your site's base URL | Yes |
| `authorName` | Author name for feed items | No |
| `authorEmailID` | Author email | No |
| `rssOutputName` | Output filename (default: `feed.xml`) | No |

## Mark Sections for RSS

In your section's `_index.md`, add `rss: true`:

```markdown
---
rss: true
---
```

This includes all posts in that section in the RSS feed.

## Link to Feed

Add the RSS link in your templates:

```html
<link rel="alternate" type="application/rss+xml"
      title="RSS Feed" href="/feed.xml">
```

Or add a visible link:

```html
<a href="/feed.xml">RSS Feed</a>
```

## Output

After building, `feed.xml` is generated in the output root directory.

## Example Setup

### 1. config.json

```json
{
  "rss": {
    "isRssEnabled": true,
    "title": "My Tech Blog",
    "link": "https://myblog.com",
    "authorName": "John Doe",
    "rssOutputName": "feed.xml"
  }
}
```

### 2. blog/_index.md

```markdown
---
rss: true
---
```

### 3. Build

```bash
./hanamark build
```

The feed will be available at `output_html/feed.xml`.
