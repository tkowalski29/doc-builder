# doc-builder

`doc-builder` rewrites the former `build-docs.sh` workflow into a self-contained Go
CLI. It scans a project for prefixed markdown files, assembles a temporary VitePress
workspace, and produces a ready-to-ship `dist` directory without requiring any
manual steps.

## Features

- Discover markdown files whose names start with a configurable prefix (default `DOC_`).
- Merge additional markdown content that already lives inside the documentation
  directory, preserving hand-crafted guides.
- Generate a VitePress sidebar by replacing the `// SIDEBAR_ITEMS - will be replaced by build script`
  placeholder in `.vitepress/base.config.js`.
- Run `npm install` (only when needed) inside the temporary workspace and execute
  `npm run docs:build` for the chosen engine (currently VitePress).
- Copy the produced `.vitepress/dist` output back into the documentation workspace.
- Provide a `helper` subcommand that explains the complete workflow and expected
  repository layout.

## Requirements

- Go 1.22 or newer (for building the CLI).
- Node.js toolchain available on the system (`npm` must be on `PATH`).
- A documentation workspace that contains:
  - `.vitepress/base.config.js` with the sidebar placeholder comment.
  - `package.json` defining `npm run docs:build`.

The command stops with a clear error message whenever a required file is missing.

## Building the CLI

```bash
go build -o bin/doc-builder ./cmd/docbuilder
```

## Usage

Run the tool from your documentation directory (the directory that contains
`.vitepress` and `package.json`). Point `--search` to the project location you
want to scan for prefixed markdown files.

```bash
./bin/doc-builder \
  --search ../ \
  --doc-dir . \
  --prefix DOC_ \
  --engine vitepress
```

### Flags

- `--search` *(required)*: root folder where prefixed markdown files are
  discovered.
- `--doc-dir` *(default: current directory)*: documentation workspace that holds
  `.vitepress` and `package.json`.
- `--prefix` *(default: `DOC_`)*: filename prefix used to select markdown sources.
- `--engine` *(default: `vitepress`)*: documentation engine. Only `vitepress`
  is implemented today; the flag keeps the interface forward compatible.
- `--temp-dir` *(default: `temp`)*: name of the temporary build directory inside
  the documentation workspace.
- `--verbose`: prints detailed progress information.

### Helper

To see a high-level overview of the pipeline, run:

```bash
./bin/doc-builder helper
```

### Example Document

Generate a sample markdown file (default `DOC_this_is_example.md`) to use as a
starting point:

```bash
./bin/doc-builder example-doc --doc-dir .
```

Use `--prefix` if your project relies on a different filename prefix.

### Typical Workflow

1. The CLI deletes and recreates `./temp` (or the configured temporary directory).
2. Markdown files that match the prefix are copied into the temp workspace with
   slugs derived from their filenames and categories read from YAML front matter.
3. Existing markdown content inside the documentation workspace is merged so that
   curated guides remain available.
4. The sidebar is regenerated based on the collected metadata and written both to
   `temp/.vitepress/config.js` and `.vitepress/config.js`.
5. Node dependencies are installed into the temp directory (unless `node_modules`
   already exists there) and `npm run docs:build` is executed.
6. The freshly generated `temp/.vitepress/dist` directory replaces
   `.vitepress/dist` in the documentation workspace.

If any of the steps fail (for example when `package.json` or the base config
cannot be found), the command stops immediately and reports the offending path in
English.

## Migrating from the Bash Script

The original `build-docs.sh` script is no longer required. The new CLI provides the
same behaviour, adds structured error handling, and keeps the door open for
additional engines beyond VitePress.

## License

This project is released under the [MIT License](LICENSE).
