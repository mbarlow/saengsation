.PHONY: help build setup setup-udev setup-group hooks check demo demo-states status clean

BINARY := saengsation

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the saengsation binary
	go build -o $(BINARY) .

setup: build setup-group setup-udev ## Full setup (build, group, udev)
	@echo ""
	@echo "Setup complete! You may need to log out/in for group changes."
	@echo "Test with: ./$(BINARY) status"

setup-udev: ## Install udev rules
	sudo cp 99-saengsation.rules /etc/udev/rules.d/
	sudo udevadm control --reload-rules
	sudo udevadm trigger
	@echo "Udev rules installed."

setup-group: ## Create plugdev group and add current user
	@if ! getent group plugdev >/dev/null 2>&1; then sudo groupadd plugdev; fi
	@sudo usermod -aG plugdev $$(whoami)
	@echo "Added $$(whoami) to plugdev."

hooks: ## Install Claude Code hooks into ~/.claude/settings.json
	@bash scripts/install-hooks.sh

check: ## Check device access and permissions
	@bash scripts/check-deps.sh

demo: build ## Run demo animations
	@bash scripts/demo.sh

demo-states: build ## Cycle through built-in states
	@bash scripts/demo-states.sh

status: build ## Show keyboard status
	./$(BINARY) status

clean: ## Remove build artifacts
	rm -f $(BINARY)
