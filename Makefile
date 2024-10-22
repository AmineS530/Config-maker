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

all : setup finish

setup:  
# checks for oh-my-zsh
	@if [ -d "$$HOME/.oh-my-zsh" ]; then \
    	echo "${GREEN}Oh-my-zsh is already installed, continuing...${NC}" ;\
	else \
  		echo "${YELLOW}Installing Oh-my-zsh...${NC}" \
		$$(wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh) && \
		sh install.sh --unattended ;\
	fi;

# Check if the directory already exists
	@if [ -d "$(destination_dir)" ]; then \
		echo "$(YELLOW)Directory $(destination_dir) already exists. Overwriting...$(NC)"; \
		rm -rf "$(destination_dir)"; \
	else \
		echo "$(CYAN)Cloning repository into $(destination_dir)...$(NC)"; \
	fi; \
	git clone --quiet --depth=1 "$(download_link)" "$(destination_dir)" && cd "$(destination_dir)"

# Execute the setup script
	@cd $(destination_dir) && zsh setup.sh

# installs docker in rootless mode
docker:
	@echo "\033[1m\033[92mGetting docker ready for first use\nPlease Wait...\033[0m"
	@curl -fsSL https://get.docker.com/rootless 2>/dev/null | sh >/dev/null 2>&1
	@echo "\033[1m\033[92mCopy paste the following command to start docker in rootless mode:\033[0m"
	export PATH=$(HOME)/bin:$$PATH
	export DOCKER_HOST=unix://$(XDG_RUNTIME_DIR)/docker.sock
	@echo "\033[1m\033[92mDocker environment is set up for rootless mode.\033[0m"

# for zsh haters
bash:
	@echo "\033[1m\033[92mRemoving zsh from being the default shell\033[0m"
	@$$(sed -i '/SHELL=\/bin\/zsh/d' ~/.bashrc )
	@$$(sed -i '/exec \/bin\/zsh -l/d' ~/.bashrc)
	@echo "\033[1m\033[92mEnjoy Bash :))\033[0m"

# make zsh default again (W move)
zsh:
	@echo "\033[1m\033[92mEnabling zsh as default shell :)\033[0m"
	@$$(printf "SHELL=/bin/zsh\nexec /bin/zsh -l\n" >>~/.bashrc)

git:
	@echo "$(YELLOW)Running git_setup.sh...$(NC)"
	@cd $(destination_dir) && \
	zsh git_setup.sh

theme:
	@echo "$(YELLOW)Running set_theme.sh...$(NC)"
	@cd $(destination_dir) && \
	zsh set_theme.sh

background:
	@echo "$(YELLOW)Running set_background.sh...$(NC)"
	@cd $(destination_dir) && \
	zsh set_background.sh


finish:
	@if pgrep gnome-terminal- > /dev/null; then \
		echo "$(CYAN)GNOME Terminal is running.$(NC) $(YELLOW)Closing all instances...$(NC)"; \
		killall gnome-terminal-server; \
		echo "$(GREEN)All GNOME Terminal instances closed successfully.$(NC)"; \
	else \
		echo "$(RED)No GNOME Terminal instances are running.$(NC)"; \
	fi
	@zenity \
	--info \
	--text="<span size=\"x-large\">Finished.</span>\n\nPlease restart your terminal to apply.\n"\
	"In case of encountering a problem, send a PM to a.sadik on Discord.\n\nTerminal will attempt to re-open after clicking\n\n<b>Happy coding :)</b>" \
	--title="Setup is done" \
	--ok-label="yessir" && \
	gnome-terminal > /dev/null 2>&1 &