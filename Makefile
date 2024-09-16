.PHONY: check-deps buf-gen gen-api-client generate

check-deps:
	@echo "Checking dependencies..."
	@which buf > /dev/null || (echo "\nError: 'buf' is not installed.\nPlease install 'buf' on macOS with:\n  brew install bufbuild/buf/buf\n" && exit 1)
	@which openapi-generator-cli > /dev/null || (echo "\nError: 'openapi-generator-cli' is not installed.\nPlease install it on macOS with:\n  brew install openapi-generator\n" && exit 1)
	@echo "All dependencies are installed."

buf-gen: check-deps
	@echo "Running buf generate..."
	buf generate

gen-api-client: check-deps
	@echo "Generating TypeScript API client..."
	openapi-generator-cli generate \
		-i gen/openapiv2/proto/server.swagger.json \
		-g typescript-fetch \
		-o ./web/src/apiClient

generate: buf-gen gen-api-client
	@echo "Code generation completed successfully."
