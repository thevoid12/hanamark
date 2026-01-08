---
created_on: 2024-12-26
---

# Configuration

Configuration is stored in `./configurables/config.json`.

## Complete Example

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
  "servePort": "3000",
  "opengraph": {
    "enabled": true,
    "siteName": "My Awesome Blog",
    "baseUrl": "https://example.com",
    "defaultImage": "/static/images/og-default.png",
    "imageWidth": "1200",
    "imageHeight": "630",
    "twitterCard": "summary_large_image"
  }
}
```

## File Paths

| Key | Description | Example |
|-----|-------------|---------|
| `sourceMDDir` | Directory containing Markdown source files | `./configurables/source_md` |
| `destHtmlDir` | Output directory for generated HTML | `./output_html` |
| `templatePath` | Directory containing HTML templates | `./configurables/templates` |
| `sourceStaticFiles` | Static files (CSS, JS, images) to copy | `./configurables/static` |
| `mdAssetsSourcePath` | Assets referenced in Markdown files | `./configurables/source_md/assets/` |
| `mdAssetsDestPath` | Destination for Markdown assets | `./output_html/assets/` |

## Index Homepage

Controls what content appears on the root `index.html`:

```json
"indexHomepageHtml": {
  "type": "section",
  "name": "blog",
  "titletag": "My Custom Site Title"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Either `section` or `page` |
| `name` | string | Name of the section or page to use |
| `titletag` | string (optional) | Custom text for the `<title>` tag of the root index.html. If not provided, defaults to the section/folder name |

| Type | Behavior |
|------|----------|
| `section` | Uses the specified section's list as the homepage |
| `page` | Copies the specified page as `index.html` |

### Examples

**Use blog section as homepage:**

```json
"indexHomepageHtml": {
  "type": "section",
  "name": "blog"
}
```

This configuration:
- Uses `list.html` template
- Generates `blog/index.html` with list of all posts
- Copies that as the root `index.html`

**With custom title tag:**

```json
"indexHomepageHtml": {
  "type": "section",
  "name": "blog",
  "titletag": "My Awesome Blog - Home"
}
```

This will use "My Awesome Blog - Home" in the `<title>` tag of the root index.html instead of the default "blog"

**Use about.html as homepage:**

```json
"indexHomepageHtml": {
  "type": "page",
  "name": "about.html"
}
```

### Nested Folder Example

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

### How List Pages Work

```
blog/
├── _index.md      -> Uses list.html   -> blog/index.html (shows links to all posts)
├── post-1.md      -> Uses single.html -> blog/post-1.html
└── post-2.md      -> Uses single.html -> blog/post-2.html
```

> **Important:** If there is no `_index.md` in a folder, no list page (`index.html`) will be generated. The engine will treat files as individual pages only.

### Custom Homepage Template

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

## OpenGraph Support (SE0)

Configure OpenGraph and Twitter Card metadata for better social media sharing.

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

## Logger

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

## Other Options

| Key | Type | Description |
|-----|------|-------------|
| `tags` | boolean | Enable tag system |
| `sortFilesByCreatedOn` | boolean | Sort lists by `created_on` date (true) or `updated` date (false) |
| `servePort` | string | Port for local development server (default: `3000`) |
