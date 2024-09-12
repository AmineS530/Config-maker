# this makefile made to set up docker, type 'make' in same directory
all:
	@echo "\033[1m\033[92mGetting docker ready for first use\033[0m"
	@curl -fsSL https://get.docker.com/rootless | sh && export PATH=$$HOME/bin:$$PATH && export DOCKER_HOST=unix://$$XDG_RUNTIME_DIR/docker.sock 
	@echo "Docker environment is set up for rootless mode."
