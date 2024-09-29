RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
CYAN=\033[0;36m
NC=\033[0m
# Clones the repo files
destination_dir=$(HOME)/Zone01_Desk_cfg

# after forking and adding ur custom background / adjusting the script to your preference put the link of your repo here (your repo should be public)
download_link=https://github.com/AmineS530/Config-maker.git


#*  <-------------------------------  Made by asadik with chatGPT and love  ------------------------------->  *#


all : setup
	@printf "Please restart your terminal to apply\nIn case of encountring a problem send a PM to a.sadik on discord"


setup: getFiles
	@cd $(destination_dir) && zsh setup.sh


getFiles:
# Check if the directory already exists

	@if [ -d "$(destination_dir)" ]; then \
		echo "$(YELLOW)Directory $(destination_dir) already exists. Overwriting...$(NC)"; \
		rm -rf "$(destination_dir)"; \
	else \
		echo "$(CYAN)Cloning repository into $(destination_dir)...$(NC)"; \
	fi; \
	git clone --quiet --depth=1 "$(download_link)" "$(destination_dir)" && cd "$(destination_dir)"

docker:
	@echo "\033[1m\033[92mGetting docker ready for first use\033[0m"
	@curl -fsSL https://get.docker.com/rootless 2>/dev/null | sh >/dev/null 2>&1
	@$$(sleep 2)
	$$(export PATH=$(HOME)/bin:$$PATH)
	$$(export DOCKER_HOST=unix://$(XDG_RUNTIME_DIR)/docker.sock)
	@echo "\033[1m\033[92mDocker environment is set up for rootless mode.\033[0m"
