IMAGE_NAME=tsql-cli
CONTAINER_NAME=tsql-cli
DOCKERFILE_PATH=./Dockerfile

build:
	docker build -t $(IMAGE_NAME) -f $(DOCKERFILE_PATH)

run:
	docker run --name $(CONTAINER_NAME) -p 3307:3306 -d $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME) && docker container rm $(CONTAINER_NAME)

kill:
	lsof -i ":3307" | awk 'NR>1 {print $$2}' | xargs kill

logs:
	docker logs $(CONTAINER_NAME)