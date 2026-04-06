.PHONY: help install setup setup-udev setup-group check demo status clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install Python dependencies
	@bash scripts/install.sh

setup: ## Full setup (group, udev, deps)
	@bash scripts/setup.sh

setup-udev: ## Install udev rules only
	sudo cp 99-saengsation.rules /etc/udev/rules.d/
	sudo udevadm control --reload-rules
	sudo udevadm trigger
	@echo "Udev rules installed."

setup-group: ## Create plugdev group and add current user
	@if ! getent group plugdev >/dev/null 2>&1; then sudo groupadd plugdev; fi
	sudo usermod -aG plugdev $$(whoami)
	@echo "Added $$(whoami) to plugdev. Log out/in or run: newgrp plugdev"

check: ## Check dependencies and device access
	@bash scripts/check-deps.sh

demo: ## Run demo animations
	@bash scripts/demo.sh

status: ## Show device connection status
	@uv run python -m saengsation status

clean: ## Remove build artifacts
	rm -rf .venv __pycache__ saengsation/__pycache__ *.egg-info dist build
