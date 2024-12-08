SERVICE_NAME 	= pismo-transactions
IMPORT_PATH 	= github.com/sathishs-dev/${SERVICE_NAME}
PKG_SRC 		= ${IMPORT_PATH}/cmd/${SERVICE_NAME}

BUILD_TAG 		= build-local
IMAGE 			= ${SERVICE_NAME}:${BUILD_TAG}
MIGRATOR_IMAGE  = pismo-migrator:${BUILD_TAG}

DOCKER_COMPOSE 	= docker compose --file docker-compose.yml

export MIGRATOR_IMAGE
export IMAGE

## Service and Migrator
build: build-service build-migrator build-unit-test

build-migrator:
	@echo "==> Building the migrator image ..."
	docker build --tag ${MIGRATOR_IMAGE} \
	 -f migrator.Dockerfile .

build-service:
	@echo "==> Building the service ..."
	docker build --pull --tag ${IMAGE} \
		--build-arg importPath=${IMPORT_PATH} \
		--build-arg pkg=${PKG_SRC} .

run-service: unit-test build-migrator build-service
	@echo "==> Launching the Service..."
	${DOCKER_COMPOSE} up -d

# Unit Test
build-unit-test:
	@echo "==> Building unit-test images..."
	docker build --file unittest.Dockerfile \
		--tag ${SERVICE_NAME}-unittest \
		--build-arg importPath=${IMPORT_PATH} .

unit-test:
	@echo "==> Runnig unit-test cases..."
	docker run --rm ${SERVICE_NAME}-unittest go test -mod vendor -v -cover -race ./...

# Utilities
mock:
	@echo "==> Generating mocks..."
	go generate ./...

down:
	${DOCKER_COMPOSE} down -v --remove-orphans

logs:
	${DOCKER_COMPOSE} logs -f

dep:
	@echo "==> Updating dependencies..."
	go mod tidy
	go mod vendor
