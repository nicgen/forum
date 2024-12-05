all: build run

name = forum

build:
	@sudo touch forum.db && docker build -t $(name) -f Dockerfile .

create:
	@if [ ! "$(shell docker ps -a -q -f name=$(name))" ]; then \
		sudo docker create --name $(name) -p 8080:8080 -v $(shell pwd)/forum.db:/app/forum.db --env-file .env $(name) && \
		echo "'$(name)' container created successfully"; \
	else \
		echo "'$(name)' container already exists"; \
	fi

start:
	@sudo docker container start $(name)

run: create start

stop:
	@sudo docker container stop $(name) && echo "$(name) container stopped successfully" || echo "No running container to stop"
# @sudo docker container stop $(name)
# @echo "$(name) container stopped successfully"

delete:
	@sudo docker container rm -f $(name) &> /dev/null && echo "$(name) container deleted correctly" || echo "No container to remove"
	@sudo docker rmi -f $(name) > /dev/null 2>&1 && echo "$(name) image deleted correctly" || echo "No image to remove"

clean: stop delete
