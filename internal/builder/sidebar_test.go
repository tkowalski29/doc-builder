package builder

import (
	"strings"
	"testing"
)

func TestBuildSectionsCreatesHierarchies(t *testing.T) {
	records := []menuRecord{
		{CategoryPath: "", Slug: "build-readme", Title: "System Overview"},
		{CategoryPath: "guides", Slug: "index", Title: "Guides Home"},
		{CategoryPath: "guides/advanced", Slug: "index", Title: "Advanced Home"},
		{CategoryPath: "guides/advanced", Slug: "deep-dive", Title: "Deep Dive"},
	}

	sections := buildSections(records)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}

	general := sections[0]
	if general.Title != "General" {
		t.Fatalf("expected first section to be 'General', got %q", general.Title)
	}
	if len(general.Items) != 1 || general.Items[0].Slug != "build-readme" {
		t.Fatalf("expected general section to contain build-readme item, got %+v", general.Items)
	}

	guides := sections[1]
	if guides.Title != "Guides" {
		t.Fatalf("expected second section title 'Guides', got %q", guides.Title)
	}
	if len(guides.Items) != 1 || guides.Items[0].Slug != "index" {
		t.Fatalf("expected guides root item to be index, got %+v", guides.Items)
	}
	if len(guides.OrderedSubs) != 1 {
		t.Fatalf("expected one subsection under guides, got %d", len(guides.OrderedSubs))
	}
	sub := guides.OrderedSubs[0]
	if sub.Title != "Advanced" {
		t.Fatalf("expected subsection title 'Advanced', got %q", sub.Title)
	}
	if len(sub.Items) != 2 {
		t.Fatalf("expected two items in subsection, got %d", len(sub.Items))
	}
	if sub.Items[0].Slug != "index" {
		t.Fatalf("expected index to be first within subsection, got %q", sub.Items[0].Slug)
	}
}

func TestRenderSidebarContainsExpectedEntries(t *testing.T) {
	records := []menuRecord{
		{CategoryPath: "platform", Slug: "tour", Title: "Platform Tour"},
		{CategoryPath: "platform/mobile", Slug: "index", Title: "Mobile"},
	}

	sections := buildSections(records)
	sidebar := renderSidebar(sections)

	checks := []string{
		"text: 'Platform'",
		"link: '/platform/tour'",
		"text: 'Mobile'",
		"link: '/platform/mobile/index'",
	}
	for _, expect := range checks {
		if !strings.Contains(sidebar, expect) {
			t.Fatalf("expected sidebar output to contain %q, got %s", expect, sidebar)
		}
	}
}
