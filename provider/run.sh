# Prebuild image
APP_NAME=tarzan_pc docker-compose build

go run main.go -owner=$1
