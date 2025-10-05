package builder

import "testing"

func TestParseFrontMatter(t *testing.T) {
	content := []byte("---\nTitle: Getting Started\nCategory: guides/intro\n---\n# Body\n")
	fm := parseFrontMatter(content)

	if got := fm["title"]; got != "Getting Started" {
		t.Fatalf("expected title to be 'Getting Started', got %q", got)
	}
	if got := fm["category"]; got != "guides/intro" {
		t.Fatalf("expected category to be 'guides/intro', got %q", got)
	}
}

func TestParseFrontMatterWithoutBlock(t *testing.T) {
	content := []byte("# Heading only\n")
	if fm := parseFrontMatter(content); len(fm) != 0 {
		t.Fatalf("expected empty map, got %v", fm)
	}
}

func TestDeriveTitleOrdering(t *testing.T) {
	content := []byte("---\ntitle: Front Matter Title\n---\n# Markdown Heading\n")
	title := deriveTitle(content, "Front Matter Title", "doc-title", "DOC_")
	if title != "Front Matter Title" {
		t.Fatalf("expected front matter title preference, got %q", title)
	}

	content = []byte("# Only Heading\n")
	title = deriveTitle(content, "", "DOC_ONLY_HEADING", "DOC_")
	if title != "Only Heading" {
		t.Fatalf("expected heading derived title, got %q", title)
	}

	title = deriveTitle([]byte("no headings here"), "", "doc_fallback", "DOC_")
	if title != "Doc Fallback" {
		t.Fatalf("expected fallback title 'Doc Fallback', got %q", title)
	}
}

func TestBuildSlugAndNormalizeCategory(t *testing.T) {
	slug := buildSlug("DOC_Example_Page.md", "DOC_")
	if slug != "example-page" {
		t.Fatalf("expected slug 'example-page', got %q", slug)
	}

	category := normalizeCategoryPath("/guides/advanced/")
	if category != "guides/advanced" {
		t.Fatalf("expected normalized category 'guides/advanced', got %q", category)
	}

	if path := normalizeCategoryPath("."); path != "" {
		t.Fatalf("expected '.' to normalize to empty string, got %q", path)
	}
}

func TestMenuKeyIncludesCategoryAndSlug(t *testing.T) {
	if key := menuKey("guides", "intro"); key != "guides|intro" {
		t.Fatalf("unexpected key %q", key)
	}
}
