local_resource(
    'firestore_emulator',
    serve_cmd='gcloud emulators firestore start --host-port=localhost:8081 --quiet',
    allow_parallel=True,
)

local_resource(
    'populate_database',
    cmd='go run server/cmd/populate_test_data/main.go',
    env={'FIRESTORE_EMULATOR_HOST': 'localhost:8081'},
    deps=['server/cmd/populate_test_data/main.go'],
    resource_deps=['firestore_emulator'],
)

local_resource(
    'run_application',
    serve_cmd='go run server/cmd/server/main.go',
    serve_env={'FIRESTORE_EMULATOR_HOST': 'localhost:8081'},
    deps=['server/cmd/server/main.go', 'server/internal/server/server.go'],
    resource_deps=['populate_database'],
)


