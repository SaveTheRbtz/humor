local_resource(
    'firestore_emulator',
    serve_cmd='gcloud emulators firestore start --host-port=localhost:8081',
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
    'leaderboard',
    cmd='bazel run //scripts/leaderboard:leaderboard_bin',
    env={'FIRESTORE_EMULATOR_HOST': 'localhost:8081'},
    deps=['scripts/leaderboard/'],
)

local_resource(
    'web',
    serve_cmd='npm install && npm start',
    serve_dir='web',
    serve_env={'REACT_APP_API_BASE_URL': 'http://localhost:8080'},
    deps=['web/src/'],
)

local_resource(
    'server',
    serve_cmd='go run server/cmd/server/main.go',
    serve_env={'FIRESTORE_EMULATOR_HOST': 'localhost:8081'},
    deps=['server/'],
    resource_deps=['populate_database'],
)
