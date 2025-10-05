package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tkowalski29/doc-builder/internal/builder"
)

func main() {
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "helper":
			runHelper()
			return
		case "example-doc":
			if err := runExampleDoc(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "unable to create example document: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	cfg := builder.Config{}
	fs := flag.NewFlagSet("doc-builder", flag.ExitOnError)
	fs.StringVar(&cfg.Prefix, "prefix", "DOC_", "File name prefix to detect documentation sources")
	fs.StringVar(&cfg.Engine, "engine", "vitepress", "Documentation engine to use (currently only 'vitepress')")
	fs.StringVar(&cfg.SearchPath, "search", "", "Root path where prefixed markdown files will be discovered")
	fs.StringVar(&cfg.DocDir, "doc-dir", ".", "Documentation workspace directory that contains .vitepress setup")
	fs.StringVar(&cfg.TempDirName, "temp-dir", "temp", "Name of the temporary build directory inside the documentation workspace")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging output")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "doc-builder rewrites prefixed markdown into a VitePress site.\n\n")
		fmt.Fprintf(fs.Output(), "Usage: doc-builder [flags]\n\n")
		fmt.Fprintf(fs.Output(), "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nRun 'doc-builder helper' to see the high-level workflow or 'doc-builder example-doc' to generate a sample markdown file.\n")
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags: %v\n", err)
		os.Exit(2)
	}

	if cfg.SearchPath == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: --search")
		fs.Usage()
		os.Exit(2)
	}

	if err := builder.New(cfg).Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
}

func runHelper() {
	message := `doc-builder helper

This CLI automates converting markdown files that follow a shared prefix into a
ready-to-ship VitePress site. The workflow is:
  1. Scan the provided search path for markdown files starting with the prefix
     (default DOC_), skipping common vendor directories.
  2. Copy matching files into a clean temporary workspace under the docs
     directory, derive slugs, titles, and categories from front matter.
  3. Merge any existing markdown that already lives in the documentation folder
     (for example curated guides) so that hand-crafted content is preserved.
  4. Compose the sidebar by translating the collected metadata into the
     // SIDEBAR_ITEMS placeholder in .vitepress/base.config.js.
  5. Install Node.js dependencies inside the temp directory when needed and run
     \'npm run docs:build\' using the selected engine (currently only VitePress).
  6. Copy the generated .vitepress/dist output back to the main docs workspace.

Typical usage:
  ./bin/doc-builder --search ../ --doc-dir . --prefix DOC_ --engine vitepress

The process stops with a descriptive error if expected files like
.vitepress/base.config.js or package.json cannot be found.
`
	fmt.Println(message)
}

func runExampleDoc(args []string) error {
	fs := flag.NewFlagSet("doc-builder example-doc", flag.ExitOnError)
	docDir := fs.String("doc-dir", ".", "Directory where the example markdown file will be created")
	prefix := fs.String("prefix", "DOC_", "Prefix to apply to the generated example filename")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: doc-builder example-doc [--doc-dir path] [--prefix PREFIX]\n")
	}

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	dir := *docDir
	if dir == "" {
		dir = "."
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("unable to resolve directory: %w", err)
	}
	if info, err := os.Stat(absDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("target directory not found: %s", absDir)
		}
		return fmt.Errorf("failed to access target directory: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("target path is not a directory: %s", absDir)
	}

	fileName := fmt.Sprintf("%sthis_is_example.md", *prefix)
	filePath := filepath.Join(absDir, fileName)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("example file already exists: %s", filePath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check for existing file: %w", err)
	}

	createdAt := time.Now().Format(time.RFC3339)
	content := fmt.Sprintf(`---
title: This Is Example
category: guides/example
description: Short note describing the purpose of this document.
last_updated: %s
---

# This Is Example

The front matter at the top defines metadata used by doc-builder:

- "title" is shown in navigation menus and as the page heading.
- "category" determines the folder hierarchy inside the generated site.
- "description" is optional but helps with search and previews.
- "last_updated" can be any string; ISO timestamps work well.

## Writing Content

Start the body with a level-one heading that repeats the title. Continue with
guides, notes, code snippets, and any additional sections required by your
documentation.

## Tips

- Place screenshots and assets next to the markdown file when possible.
- Keep the filename descriptive—doc-builder converts underscores to hyphens for
  links.
- You can add more keys to the front matter if the theme supports them.
`, createdAt)

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to create example file %s: %w", filePath, err)
	}

	fmt.Printf("Example markdown created at %s\n", filePath)
	return nil
}
