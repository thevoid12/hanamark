---
created_on: 2024-12-26
---

# Content Organization

Hanamark mirrors your source directory structure in the output.

## Directory Mapping

```
source_md/                    output_html/
├── about.md            ->    ├── about.html
├── projects.md         ->    ├── projects.html
└── blog/                     └── blog/
    ├── _index.md       ->        ├── index.html (list page)
    ├── post-1.md       ->        ├── post-1.html
    └── post-2.md       ->        └── post-2.html
```

## Single Pages

Any `.md` file (except `_index.md`) becomes a single HTML page:

```markdown
---
created_on: 2024-01-15
tags: ["about", "personal"]
---

# About Me

This is my about page content...
```

**Output:** `about.html`

## List Pages (Sections)

Directories with a `_index.md` file become **sections** with automatic list pages.

> **Important:** If you want a list page, it is mandatory to have `_index.md` in that folder.

```
blog/
├── _index.md           # Section configuration
├── post-1.md           # Blog post
├── post-2.md           # Blog post
└── nested/             # Nested section
    ├── _index.md
    └── deep-post.md
```

The `_index.md` file configures the section:

```markdown
---
rss: true
template: "blog/custom_list_template.html"
---
```

**Output:** `blog/index.html` with links to all posts in the section.

## Nested Sections

You can nest sections as deep as you need:

```
docs/
├── _index.md
├── intro.md
└── guides/
    ├── _index.md
    ├── beginner.md
    └── advanced.md
```

Each section with `_index.md` gets its own list page.
