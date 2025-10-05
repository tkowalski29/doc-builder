package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type section struct {
	Key         string
	Title       string
	Items       []menuRecord
	subsections map[string]*subsection
	OrderedSubs []*subsection
}

type subsection struct {
	Key   string
	Title string
	Items []menuRecord
}

func (b *Builder) generateConfig(env environment, records []menuRecord) error {
	if b.cfg.Verbose {
		fmt.Println("[4/7] Generating VitePress sidebar configuration")
	}

	sections := buildSections(records)
	sidebar := renderSidebar(sections)

	baseData, err := os.ReadFile(env.baseConfig)
	if err != nil {
		return fmt.Errorf("failed to read base config %s: %w", env.baseConfig, err)
	}

	const placeholder = "// SIDEBAR_ITEMS - will be replaced by build script"
	var output []byte

	baseString := string(baseData)
	if strings.Contains(baseString, placeholder) {
		output = []byte(strings.Replace(baseString, placeholder, sidebar, 1))
	} else if len(sections) > 0 {
		return fmt.Errorf("placeholder '%s' not found in %s", placeholder, env.baseConfig)
	} else {
		output = baseData
	}

	tempConfig := filepath.Join(env.tempDir, ".vitepress", "config.js")
	if err := os.MkdirAll(filepath.Dir(tempConfig), 0o755); err != nil {
		return fmt.Errorf("failed to create temp config directory: %w", err)
	}
	if err := os.WriteFile(tempConfig, output, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", tempConfig, err)
	}
	if err := os.WriteFile(env.outputConfig, output, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", env.outputConfig, err)
	}

	return nil
}

func buildSections(records []menuRecord) []*section {
	sectionsMap := map[string]*section{}
	for _, rec := range records {
		category := normalizeCategoryPath(rec.CategoryPath)
		parts := []string{}
		if category != "" {
			parts = strings.Split(category, "/")
		}

		rootKey := ""
		if len(parts) > 0 {
			rootKey = parts[0]
		}

		sec, exists := sectionsMap[rootKey]
		if !exists {
			title := "General"
			if rootKey != "" {
				title = formatTitle(rootKey)
			}
			sec = &section{
				Key:         rootKey,
				Title:       title,
				Items:       []menuRecord{},
				subsections: map[string]*subsection{},
			}
			sectionsMap[rootKey] = sec
		}

		if len(parts) <= 1 {
			sec.Items = append(sec.Items, rec)
			continue
		}

		subKey := strings.Join([]string{rootKey, parts[1]}, "/")
		sub, exists := sec.subsections[subKey]
		if !exists {
			sub = &subsection{
				Key:   subKey,
				Title: formatTitle(parts[1]),
				Items: []menuRecord{},
			}
			sec.subsections[subKey] = sub
		}
		sub.Items = append(sub.Items, rec)
	}

	keys := make([]string, 0, len(sectionsMap))
	for key := range sectionsMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ordered := make([]*section, 0, len(keys))
	for _, key := range keys {
		sec := sectionsMap[key]
		sortMenuRecords(sec.Items)
		subKeys := make([]string, 0, len(sec.subsections))
		for subKey := range sec.subsections {
			subKeys = append(subKeys, subKey)
		}
		sort.Strings(subKeys)
		orderedSubs := make([]*subsection, 0, len(subKeys))
		for _, subKey := range subKeys {
			sub := sec.subsections[subKey]
			sortMenuRecords(sub.Items)
			orderedSubs = append(orderedSubs, sub)
		}
		sec.OrderedSubs = orderedSubs
		ordered = append(ordered, sec)
	}

	return ordered
}

func sortMenuRecords(items []menuRecord) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Slug == "index" && items[j].Slug != "index" {
			return true
		}
		if items[i].Slug != "index" && items[j].Slug == "index" {
			return false
		}
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})
}

func renderSidebar(sections []*section) string {
	if len(sections) == 0 {
		return ""
	}

	var lines []string
	for _, sec := range sections {
		lines = append(lines, "      {")
		lines = append(lines, fmt.Sprintf("        text: '%s',", escapeQuotes(sec.Title)))
		lines = append(lines, "        collapsed: true,")
		lines = append(lines, "        items: [")
		for _, item := range sec.Items {
			lines = append(lines, renderSidebarItem(item, 2))
		}
		for _, sub := range sec.OrderedSubs {
			lines = append(lines, "          {")
			lines = append(lines, fmt.Sprintf("            text: '%s',", escapeQuotes(sub.Title)))
			lines = append(lines, "            collapsed: true,")
			lines = append(lines, "            items: [")
			for _, item := range sub.Items {
				lines = append(lines, renderSidebarItem(item, 3))
			}
			lines = append(lines, "            ]")
			lines = append(lines, "          },")
		}
		lines = append(lines, "        ]")
		lines = append(lines, "      },")
	}

	return strings.Join(lines, "\n")
}

func renderSidebarItem(item menuRecord, depth int) string {
	indent := strings.Repeat("  ", depth)
	linkParts := []string{item.CategoryPath, item.Slug}
	link := strings.Trim(strings.Join(linkParts, "/"), "/")
	return fmt.Sprintf("%s{ text: '%s', link: '/%s' },", indent, escapeQuotes(item.Title), link)
}
