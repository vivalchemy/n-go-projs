.PHONY: gen/proto gen

help:
	@bat --language=Makefile --paging=never Makefile

gen: gen/proto

gen/proto:
	@echo "Generating proto files"
	@protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

dev:
	@echo "Starting development containers"
	@docker compose --profile development up --build --watch

prod:
	@echo "Starting production containers"
	@docker compose --profile production up --build

down:
	@echo "Stopping containers"
	@docker compose down

clean:
	@echo "Cleaning proto files"
	@rm -rf pb/*.pb.go
	@echo "Cleaning tmp directory"
	@rm -rf tmp
