.PHONY: check-deps buf-gen gen-api-client generate build deploy

check-deps:
	@echo "Checking dependencies..."
	@which buf > /dev/null || (echo "\nError: 'buf' is not installed.\nPlease install 'buf' on macOS with:\n  brew install bufbuild/buf/buf\n" && exit 1)
	@which openapi-generator-cli > /dev/null || (echo "\nError: 'openapi-generator-cli' is not installed.\nPlease install it on macOS with:\n  brew install openapi-generator\n" && exit 1)
	@echo "All dependencies are installed."

buf-gen: check-deps
	@echo "Running buf generate..."
	buf generate --path proto/

gen-api-client: check-deps
	@echo "Generating TypeScript API client..."
	openapi-generator-cli generate \
		-i gen/openapiv2/proto/server.swagger.json \
		-g typescript-fetch \
		-o ./web/src/apiClient

generate: buf-gen gen-api-client
	@echo "Code generation completed successfully."

build:
	@echo "Building server Docker images..."
	docker build -f arena.Dockerfile --platform linux/amd64 -t gcr.io/humor-arena/server:latest --target app .
	docker build -f leaderboard.Dockerfile --platform linux/amd64 -t gcr.io/humor-arena/leaderboard:latest --target app .

push:
	@echo "Pushing Docker images to Google Container Registry..."
	docker push gcr.io/humor-arena/server:latest
	docker push gcr.io/humor-arena/leaderboard:latest

deploy: build push
	@echo "Deploying to Google Cloud Run..."
	gcloud run deploy humor-arena \
		--image gcr.io/humor-arena/server:latest \
		--platform managed \
		--region us-central1 \
		--allow-unauthenticated \
		--service-account cloud-run-firestore-sa@humor-arena.iam.gserviceaccount.com
	gcloud run jobs update leaderboard \
		--image gcr.io/humor-arena/leaderboard:latest \
		--task-timeout 10m \
		--region us-central1 \
		--service-account cloud-run-firestore-sa@humor-arena.iam.gserviceaccount.com
	gcloud run jobs execute leaderboard --region us-central1
	@echo "Deployment completed successfully."
