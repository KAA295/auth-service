generate_docs:
	swag init -g cmd/app/main.go

build: generate_docs
	docker compose build

launch:
	docker compose up -d

build_launch: build launch

build_launch_test: build_launch
	docker compose --profile tests up -d

stop:
	docker compose down
