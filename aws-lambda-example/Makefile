help:
	@echo "gen - generate the binary for aws lambda and zip it"
	@echo "help - this help"
	@echo "gen/bin - generate the binary for linux/arm64"
	@echo "gen/zip - generate the zip for aws lambda"
	@echo "clean - remove the binary"

gen: gen/bin gen/zip

gen/bin:
	@echo "Building for Linux/ARM64"
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap ./main.go

gen/zip:
	@echo "Zipping for aws lambda"
	zip boot.zip bootstrap

clean:
	@echo "Cleaning the generated files"
	@rm -f bootstrap boot.zip
