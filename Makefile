SERVICE_NAME 	= pismo-transactions
IMPORT_PATH 	= github.com/sathishs-dev/${SERVICE_NAME}
PKG_SRC 		= ${IMPORT_PATH}/cmd/${SERVICE_NAME}

BUILD_TAG 		= build-local
IMAGE 			= ${SERVICE_NAME}:${BUILD_TAG}
MIGRATOR_IMAGE  = pismo-migrator:${BUILD_TAG}

DOCKER_COMPSOE 	= docker compose --file docker-compose.yml

export MIGRATOR_IMAGE

build:

build-migrator:
	docker build --tag ${MIGRATOR_IMAGE} \
	 -f migrator.Dockerfile .

migrate: build-migrator
	${DOCKER_COMPSOE} up --no-recreate -d 

down:
	${DOCKER_COMPSOE} down -v --remove-orphans

logs:
	${DOCKER_COMPSOE} logs -f

dep:
	go mod tidy
	go mod vendor
