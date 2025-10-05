package builder

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"
	"unicode"
)

func parseFrontMatter(content []byte) map[string]string {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	fm := map[string]string{}
	if !scanner.Scan() {
		return fm
	}
	first := strings.TrimSpace(scanner.Text())
	if first != "---" {
		return fm
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "---" {
			break
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")
		fm[key] = value
	}
	return fm
}

func deriveTitle(content []byte, frontMatterTitle string, slug string, prefix string) string {
	if frontMatterTitle != "" {
		return frontMatterTitle
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}

	cleaned := slug
	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		cleaned = "Untitled"
	}
	return toTitleCase(cleaned)
}

func buildSlug(fileName, prefix string) string {
	slug := strings.TrimSuffix(fileName, ".md")
	if prefix != "" && strings.HasPrefix(slug, prefix) {
		slug = strings.TrimPrefix(slug, prefix)
	}
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ToLower(slug)
	return slug
}

func normalizeCategoryPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	path = filepath.ToSlash(path)
	path = strings.Trim(path, "/")
	if path == "." {
		return ""
	}
	return path
}

func menuKey(category, slug string) string {
	return category + "|" + slug
}

func toTitleCase(value string) string {
	words := strings.Fields(value)
	if len(words) == 0 {
		return value
	}
	for i, word := range words {
		words[i] = capitalize(word)
	}
	return strings.Join(words, " ")
}

func capitalize(word string) string {
	if word == "" {
		return word
	}
	runes := []rune(strings.ToLower(word))
	runes[0] = unicode.ToTitle(runes[0])
	return string(runes)
}

func formatTitle(input string) string {
	replaced := strings.ReplaceAll(input, "_", " ")
	replaced = strings.ReplaceAll(replaced, "-", " ")
	replaced = strings.TrimSpace(replaced)
	if replaced == "" {
		return "Untitled"
	}
	words := strings.Fields(replaced)
	for i, word := range words {
		words[i] = capitalize(word)
	}
	return strings.Join(words, " ")
}
