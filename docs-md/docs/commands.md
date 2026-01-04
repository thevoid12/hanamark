---
created_on: 2024-12-26
---

# CLI Commands

Reference for all Hanamark command-line commands.

## hanamark init

Initialize a new project in the current directory.

```bash
./hanamark init
```

Creates a `configurables/` directory with:
- Default `config.json`
- Sample templates
- Example content
- Static file structure

## hanamark build

Build the static site.

```bash
./hanamark build
```

This command:
- Parses all Markdown files
- Applies templates
- Copies static assets
- Generates tag pages (if enabled)
- Creates RSS feed (if enabled)
- Outputs to `destHtmlDir`

## hanamark serve

Start a local development server.

```bash
./hanamark serve
```

Serves the built site at `http://localhost:3000` (or configured port).

> **Note:** Run `hanamark build` first to generate the site.

## hanamark help

Display help information.

```bash
./hanamark help
```

## hanamark -version

Show version information.

```bash
./hanamark -version
```

## Typical Workflow

```bash
# 1. Initialize project
./hanamark init

# 2. Edit content in configurables/source_md/

# 3. Build the site
./hanamark build

# 4. Preview locally
./hanamark serve

# 5. Deploy output_html/ to your hosting provider
```

## Configuration

Commands use settings from `./configurables/config.json`. Key settings:

| Setting | Command | Description |
|---------|---------|-------------|
| `destHtmlDir` | build | Output directory |
| `servePort` | serve | Server port (default: 3000) |
| `logger.level` | all | Logging verbosity |
