APP := photosort
SRC := ./cmd/$(APP)
OUT := bin

LDFLAGS := -s -w

PLATFORMS := \
	darwin/arm64 \
	windows/amd64

.PHONY: clean build-all
.DEFAULT_GOAL := build-all

clean:
	rm -rf $(OUT)

build-all: clean
	@mkdir -p $(OUT)
	@echo "==> building $(APP) for: $(PLATFORMS)"
	@for p in $(PLATFORMS); do \
		GOOS=$${p%/*}; GOARCH=$${p#*/}; \
		EXT=""; [ "$$GOOS" = "windows" ] && EXT=".exe"; \
		OUTFILE="$(OUT)/$(APP)-$$GOOS-$$GOARCH$$EXT"; \
		echo "  -> $$OUTFILE"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -trimpath -ldflags "$(LDFLAGS)" -o $$OUTFILE $(SRC); \
	done
	@echo "==> done. outputs in ./$(OUT)"
