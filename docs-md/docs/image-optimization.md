---
created_on: 2026-07-02
---

# Image Optimization

Hanamark automatically converts images referenced in Markdown into responsive, cache-friendly `<picture>` markup at build time — resizing to multiple widths, encoding to WebP (with a same-format raster fallback for browsers that don't support WebP), and emitting explicit `width`/`height` attributes so pages don't suffer layout shift while images load.

## How It Works

```markdown
![A photo](./assets/photo.png)
```

becomes, at build time:

```html
<picture>
  <source type="image/webp" srcset="/assets/generated/photo-...-400w.webp 400w, /assets/generated/photo-...-800w.webp 800w" sizes="(max-width: 800px) 100vw, 800px">
  <img src="/assets/generated/photo-...-800w.png" width="800" height="533" alt="A photo" loading="lazy" decoding="async" srcset="..." sizes="...">
</picture>
```

No extra markup or manual resizing is needed — every image referenced in Markdown goes through this pipeline automatically as long as image processing is enabled (the default).

## Configuration

```json
{
  "image": {
    "enabled": true,
    "outputFormat": "webp",
    "backupOriginalFormat": true,
    "quality": {
      "webp": 82
    },
    "presets": {
      "banner": {
        "width": 1600,
        "breakpoints": [480, 800, 1200, 1600],
        "loading": "eager",
        "fetchpriority": "high",
        "sizes": "100vw"
      },
      "content": {
        "width": 800,
        "breakpoints": [400, 800],
        "loading": "lazy",
        "fetchpriority": "auto",
        "sizes": "(max-width: 800px) 100vw, 800px"
      }
    }
  }
}
```

| Key | Description | Default |
|-----|-------------|---------|
| `enabled` | Turn image processing on/off | `true` |
| `outputFormat` | Modern format every image is converted to | `webp` |
| `backupOriginalFormat` | Also generate a same-format-as-source raster fallback (jpg/png) for browsers without WebP support | `true` |
| `quality.webp` | WebP encode quality (1-100) | `82` |
| `presets` | Named sets of sizing/HTML-attribute rules — see below | — |

> **Note:** Setting `enabled: false` restores gomarkdown's plain `<img>` output with no resizing, caching, or extra markup.

## Presets

A preset controls how big an image is generated and which HTML attributes it gets. Hanamark's default config ships two preset names, but any name works as long as it's defined under `image.presets`:

| Preset | Typical use | Loading | Fetch Priority |
|--------|-------------|---------|-----------------|
| `banner` | Hero/first image on a page | `eager` | `high` |
| `content` | Every other image | `lazy` | `auto` |

### Preset Fields

| Field | Description |
|-------|-------------|
| `width` | Target width (px) images are resized to for this preset |
| `breakpoints` | Additional smaller widths generated for `srcset`, so browsers only download what they need |
| `loading` | `loading=""` attribute (`eager` or `lazy`) |
| `fetchpriority` | `fetchpriority=""` attribute (`high`, `low`, or `auto`) |
| `sizes` | `sizes=""` attribute, used by the browser to pick the right `srcset` candidate |

Every image not explicitly annotated uses the `content` preset. A page's very first image can opt into the `banner` preset via front matter:

```markdown
---
created_on: 2026-01-01
first_image_preset: banner
---
```

See [Front Matter](/docs/front-matter.html) for details.

## Per-Image Directives

Individual images can override the preset, or just a single attribute, without touching front matter — by appending a query string to the image's path in Markdown:

```markdown
![Hero image](./assets/hero.jpg?preset=banner)
![Custom width](./assets/photo.png?w=1000)
![Custom height](./assets/portrait.png?h=1200)
![High-priority content image](./assets/important.png?fetchpriority=high)
```

| Param | Description |
|-------|-------------|
| `preset` | Use a specific named preset for this image only |
| `w` | Explicit target width in pixels (overrides the preset's width) |
| `h` | Explicit target height in pixels — width is derived to preserve aspect ratio. Ignored if `w` is also given |
| `fetchpriority` | Overrides just the `fetchpriority=""` attribute for this image (`high`, `low`, or `auto`) — `loading` and every other attribute still come from the preset |

Directives can be combined, e.g. `?w=1000&fetchpriority=high`. This is useful for marking a specific mid-page image as high-priority (for example, one Lighthouse flags as part of the critical rendering path) without moving it into the `banner` preset or changing its `loading` behavior.

An invalid or unrecognized `fetchpriority` value is silently ignored and the preset's own value is used instead.

## Build Cache

Generated variants are cached under `hanamark_internal/` so unchanged images skip re-encoding on later builds. The cache key (source bytes + width + format + quality) is baked directly into each cached file's name, so there's no separate index to go stale — it's safe to delete `hanamark_internal/` at any time; it's rebuilt automatically on the next build.
