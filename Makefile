DIST    := ./dist
TUI_DIR := ./tui

TUI_BIN    := $(DIST)/alpha-aide-tui

.PHONY: all build build-tui run-tui \
        run-tui-anthropic run-tui-gemini run-tui-openai \
        clean

# ── Build ──────────────────────────────────────────────────────────────────

all: build

build: build-tui

build-tui:
	mkdir -p $(DIST)
	go build -C $(TUI_DIR) -o ../$(TUI_BIN) .

# ── Run TUI ────────────────────────────────────────────────────────────────

run-tui: build-tui
	$(TUI_BIN)

run-tui-anthropic: build-tui
	$(TUI_BIN) -provider anthropic -model claude-haiku-4-5-20251001

run-tui-gemini: build-tui
	$(TUI_BIN) -provider gemini -model gemini-2.0-flash

run-tui-openai: build-tui
	$(TUI_BIN) -provider openai -model gpt-4o-mini

# ── Misc ───────────────────────────────────────────────────────────────────

clean:
	rm -rf $(DIST)
