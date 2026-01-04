---
created_on: 2024-12-26
---

# Quick Start

Get a Hanamark site up and running in under 5 minutes.

## Step 1: Initialize Project

```bash
./hanamark init
```

This creates a `configurables/` directory with:
- Default configuration
- Sample templates
- Example content
- Static file structure

## Step 2: Add Content

Create a Markdown file in `configurables/source_md/`:

```markdown
---
created_on: 2024-12-26
tags: ["hello", "first-post"]
---

# My First Post

Welcome to my new site built with Hanamark!
```

## Step 3: Build

```bash
./hanamark build
```

This generates your static site in `output_html/`.

## Step 4: Preview

```bash
./hanamark serve
```

Open `http://localhost:3000` in your browser.

## What's Next?

- Edit `configurables/config.json` to customize your site
- Modify templates in `configurables/templates/`
- Add CSS in `configurables/static/css/`
- Create more content in `configurables/source_md/`
