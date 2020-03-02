build-local:
	go build -ldflags "-X main.commit=`git rev-parse HEAD`" -gcflags=-trimpath=$$GOPATH -o ./deployment/requester-service-local .

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.commit=`git rev-parse HEAD`" -gcflags=-trimpath=$$GOPATH -asmflags=-trimpath=$$GOPATH -o ./deployment/requester-service-linux .

docker-build-dev: build-linux
	docker build -f ./deployment/Dockerfile-dev -t gcr.io/gems-org/requester-service-dev:$(version) ./deployment
	rm ./deployment/requester-service-linux

docker-build-prod: build-linux
	docker build -f ./deployment/Dockerfile-prod -t gcr.io/gems-org/requester-service:$(version) ./deployment
	rm ./deployment/requester-service-linux

run: build-local
	cd ./deployment && ./requester-service-local --env=compose

run-local: build-local
	cd ./deployment && ./requester-service-local --env=local

get-credentials-dev:
	gcloud container clusters get-credentials dev --zone=us-central1-a

get-credentials-prod:
	gcloud container clusters get-credentials api --zone=us-central1-a

push-dev:
	gcloud docker -- push gcr.io/gems-org/requester-service-dev:$(version)

push-prod:
	gcloud docker -- push gcr.io/gems-org/requester-service:$(version)

create-dev: docker-build-dev push-dev
	kubectl create -f ./deployment/deployment.yaml
	kubectl create -f ./deployment/service.yaml

