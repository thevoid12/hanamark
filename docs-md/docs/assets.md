---
created_on: 2024-12-26
---

# Assets & Static Files

Hanamark handles two types of assets: static files and Markdown assets.

## Static Files

Files in `sourceStaticFiles` are copied to `destHtmlDir/static/`. All static assets are served from the `/static/` path.

```
configurables/static/          output_html/static/
├── css/                  ->   ├── css/
│   └── styles.css             │   └── styles.css
├── images/               ->   ├── images/
│   └── logo.png               │   └── logo.png
└── favicon.ico           ->   └── favicon.ico
```

## Standard Asset Paths

| Asset Type | Path |
|------------|------|
| CSS | `/static/css/styles.css` |
| Favicon | `/static/favicon.ico` |
| Images | `/static/images/your-image.png` |
| JavaScript | `/static/js/script.js` |
| Fonts | `/static/fonts/font.woff2` |

## Using in Templates

Reference static assets in your templates:

```html
<link rel="icon" href="/static/favicon.ico">
<link rel="stylesheet" href="/static/css/styles.css">
<img src="/static/images/logo.png" alt="Logo">
<script src="/static/js/script.js"></script>
```

> **Important:** Always use absolute paths starting with `/static/` to ensure assets load correctly on all pages regardless of URL depth.

## Markdown Assets

Images referenced in Markdown are copied from `mdAssetsSourcePath` to `mdAssetsDestPath`.

### Configuration

```json
{
  "filepath": {
    "mdAssetsSourcePath": "./configurables/source_md/assets/",
    "mdAssetsDestPath": "./output_html/assets/"
  }
}
```

### Usage in Markdown

```markdown
![My Image](./assets/photo.png)
```

The image will be copied to `output_html/assets/photo.png` and the path will be updated in the generated HTML.

## Directory Structure

```
configurables/
├── source_md/
│   ├── assets/           # Markdown assets
│   │   ├── photo.png
│   │   └── diagram.svg
│   └── blog/
│       └── post.md       # References ./assets/photo.png
└── static/               # Static files
    ├── css/
    ├── images/
    └── favicon.ico
```

## Best Practices

- Keep static files (CSS, JS, fonts) in `static/`
- Keep content images in `source_md/assets/`
- Use descriptive filenames
- Optimize images before adding them
