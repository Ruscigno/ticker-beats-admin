# set the IMAGE_NAME variable
# It should consist of repo name, image name and version
# use hyphens as separator
IMAGE_NAME = gcr.io/ticker-beats/ticker-beats-admin
# set the IMAGE_VERSION variable
# It should consist of the version of the image
# use hyphens as separator
IMAGE_VERSION = 1
# set the IMAGE_TAG variable
# It should consist of the IMAGE_NAME and IMAGE_VERSION
# use colon as separator
IMAGE_TAG = $(IMAGE_NAME):$(IMAGE_VERSION)
# set the IMAGE_LATEST variable
# It should consist of the IMAGE_NAME and latest
# use colon as separator
IMAGE_LATEST = $(IMAGE_NAME):latest
IMAGE_PREVIOUS = $(IMAGE_NAME):previous


.PHONY: generate
generate: ## Traverses project recursively, running go generate commands
	$(GO) generate ./...

.PHONY: build
build: ## builds the project using docker build
	docker build -t $(IMAGE_TAG) .
	docker tag $(IMAGE_LATEST) $(IMAGE_PREVIOUS)
	docker tag $(IMAGE_TAG) $(IMAGE_LATEST)
	docker push $(IMAGE_TAG)
	docker push $(IMAGE_LATEST)

.PHONY: pg_dump
## starts a postres 15 docker container and dumps the 'tickerheart'@super-server database to a file
pg_dump:
	docker rm -f pg-container-pg_dump
	docker run --name pg-container-pg_dump -e POSTGRES_PASSWORD=mypwd -d postgres:15
	echo "PGPASSWORD=iApFDLaGymdFzvneMaYVRoLJqYSCgCHU pg_dump -U ticker-beats -d tickerheart -h super-server -p 50432 > tickerheart.sql" | docker exec -i pg-container-pg_dump sh -c 'cat > backup.sh'
	docker exec pg-container-pg_dump chmod +x backup.sh
	docker exec pg-container-pg_dump /bin/sh -c /backup.sh
	docker cp pg-container-pg_dump:/tickerheart.sql .
	docker rm -f pg-container-pg_dump
	PGPASSWORD=mypassword psql -U ticker-beats -h localhost -p 5432 -d metatrader5 < tickerheart.sql

	# PGPASSWORD=mypassword pg_dump -U ticker-beats -d metatrader5 -h localhost -p 5432 > metatrader5.sql
	# PGPASSWORD=mypassword psql -U ticker-beats -h localhost -p 5432 -d metatrader5 < metatrader5.sql