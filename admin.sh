dlv attach 13654 --headless --listen=:2345 --log --api-version=2 --accept-multiclient --continue
dlv attach 63481 --headless --listen=:2346 --log --api-version=2 --accept-multiclient --continue

swag init -g cmd/openim-api/main.go -o cmd/swagger/docs --parseDependency --parseInternal
swag init -g cmd/api/admin-api/main.go -o cmd/api/admin-api/docs --parseDependency --parseInternal

dlv attach 4889 --headless --listen=:2345 --log --api-version=2 --accept-multiclient --continue