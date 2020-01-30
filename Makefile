DOCKER_IMAGE ?= dvc_discord_api_server
DOCKER_CONTAINER ?= dvc_discord_api_server
HIDE ?= @
PORT ?= 7080
HOSTPORT ?= 7080
NETWORK ?= bridge
SERVICES = $(DOCKER_IMAGE)

build:
	$(HIDE)docker-compose -f docker/docker-compose.yml build $(SERVICES)

start:
	$(HIDE)docker-compose -f docker/docker-compose.yml up --build $(DOCKER_CONTAINER)

daemon:
	$(HIDE)docker-compose -f docker/docker-compose.yml up -d --build $(DOCKER_CONTAINER)

stop: 
	$(HIDE)docker stop $(DOCKER_CONTAINER)
	$(HIDE)docker container rm $(DOCKER_CONTAINER)
