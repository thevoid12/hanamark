---
created_on: 2024-12-26
---

# Examples

Practical examples to help you get started with Hanamark.
- Check out this repo on how I created this blog using hanamark :[link](https://github.com/thevoid12/hanamark-doc-template)

## Basic Blog Post

`source_md/blog/my-first-post.md`:

```markdown
---
created_on: 2024-01-15
tags: ["introduction", "personal"]
---

# My First Post

Welcome to my blog! This is my first post using **Hanamark**.

## What I'll Write About

- Programming tutorials
- Project updates
- Random thoughts

![A nice image](./assets/hero.png)

Thanks for reading!
```

## Draft Post

Posts with `draft: true` are excluded from the build:

```markdown
---
created_on: 2024-01-20
draft: true
---

# Work in Progress

This post won't appear in the build until draft is removed or set to false.
```

## Custom Template Post

Use a specific template for special content:

```markdown
---
created_on: 2024-01-25
template: "blog/featured_post.html"
tags: ["featured"]
---

# Featured Article

This post uses a special template for featured content.
```

## Section Configuration

`source_md/blog/_index.md`:

```markdown
---
rss: true
template: "blog/custom_list_template.html"
---
```

This section will be included in RSS and uses a custom list template.

## Complete Template Example

`templates/blog/single.html`:

```html
<article class="blog-post">
  <header>
    <h1>{{ .PageTitle }}</h1>
    <div class="meta">
      <time datetime="{{ .CreatedDate.Format "2006-01-02" }}">
        {{ .CreatedDate.Format "January 2, 2006" }}
      </time>
      {{ if .Tags }}
      <div class="tags">
        {{ range .Tags }}
          <a href="{{ .TagDestPath }}" class="tag">{{ .TagName }}</a>
        {{ end }}
      </div>
      {{ end }}
    </div>
  </header>

  <div class="content">
    {{ .GenHtml }}
  </div>

  <footer>
    <p>Last updated: {{ .UpdatedDate.Format "Jan 2, 2006" }}</p>
  </footer>
</article>
```

## Complete List Template

`templates/blog/list.html`:

```html
<section class="blog-list">
  <h1>{{ .PageTitle }}</h1>
  <ul>
    {{ range .List }}
      <li class="post-item">
        <a href="{{ .DestPageDir }}">{{ .PageTitle }}</a>
        <time>{{ .CreatedDate.Format "Jan 2, 2006" }}</time>
      </li>
    {{ end }}
  </ul>
</section>
```

## Minimal config.json

```json
{
  "filepath": {
    "sourceMDDir": "./configurables/source_md",
    "destHtmlDir": "./output_html",
    "templatePath": "./configurables/templates",
    "sourceStaticFiles": "./configurables/static"
  },
  "tags": true,
  "servePort": "3000"
}
```

## Full config.json

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
  "servePort": "3000"
}
```
