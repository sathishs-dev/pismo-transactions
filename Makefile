SERVICE_NAME 	= pismo-transactions
IMPORT_PATH 	= github.com/sathishs-dev/${SERVICE_NAME}
PKG_SRC 		= ${IMPORT_PATH}/cmd/${SERVICE_NAME}

BUILD_TAG 		= build-local
IMAGE 			= ${SERVICE_NAME}:${BUILD_TAG}
MIGRATOR_IMAGE  = pismo-migrator:${BUILD_TAG}

DOCKER_COMPOSE 	= docker compose --file docker-compose.yml

export MIGRATOR_IMAGE
export IMAGE

build: build-service build-migrator

build-service:
	@echo "==> Building the service ..."
	docker build --pull --tag ${IMAGE} \
		--build-arg importPath=${IMPORT_PATH} \
		--build-arg pkg=${PKG_SRC} .

build-migrator:
	docker build --tag ${MIGRATOR_IMAGE} \
	 -f migrator.Dockerfile .

run-service: build-migrator build-service
	${DOCKER_COMPOSE} up --no-recreate -d 
	docker compose logs -f pismo-transactions

down:
	${DOCKER_COMPOSE} down -v --remove-orphans

logs:
	${DOCKER_COMPOSE} logs -f

dep:
	go mod tidy
	go mod vendor
