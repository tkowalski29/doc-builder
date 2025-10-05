.PHONY: build docs preview clean

ROOT := ../..
DOC_DIR := .doc
BINARY := doc-builder

build:
	go build -o $(BINARY) $(ROOT)/cmd/docbuilder

docs: build
	cd $(DOC_DIR) && ../$(BINARY) --search .. --doc-dir . --prefix DOC_ --engine vitepress

preview: docs
	cd $(DOC_DIR) && npx serve .vitepress/dist

clean:
	rm -f $(BINARY)
	rm -rf $(DOC_DIR)/temp $(DOC_DIR)/.vitepress/dist
