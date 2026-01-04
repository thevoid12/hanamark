---
created_on: 2024-12-26
---

# Getting Started

Hanamark is a fast, minimal, and flexible static site generator written in Go. I wrote this engine for my personal blog [thisisvoid.in](https://www.thisisvoid.in) and now putting it out so that anyone can quickly build and maintain their static sites.

## What is Hanamark?

Hanamark converts your Markdown files into a complete static website. It handles:

- **Markdown parsing** - Write content in Markdown, get HTML
- **Templating** - Use Go templates to control your site's look
- **Content organization** - Automatic list pages for sections
- **Tags** - Organize and categorize your content
- **RSS feeds** - Built-in feed generation

## Quick Start

```bash
# 1. Download the binary for your operating system or build the binary by yourself from the source code using the command
go build -o hanamark .go build -o hanamark .

# 2. Initialize a new project
./hanamark init

# 3. Build the site
./hanamark build

# 4. Preview locally
./hanamark serve
```

By default (can be configured in config.json) if served,

 Your site will be available at

`http://localhost:3000`.

## Next Steps

- [Installation](/docs/installation.html) - Download or build Hanamark
- [Project Structure](/docs/project-structure.html) - Understand the directory layout
- [Configuration](/docs/configuration.html) - Configure your site
- [Rules and Constraints](/docs/rules-and-constraints.html) - What you can and cannot do
