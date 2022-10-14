OUTPUT_DIR := ./bin
OUTPUT_API := ${OUTPUT_DIR}/api
DEBUG_PORT_API := :8181

build.api:
	go build -o ${OUTPUT_API} cmd/api/main.go
build.api.dev:
	go build -gcflags="all=-N -l" -o ${OUTPUT_API} cmd/api/main.go
run.api: clean build.api
	${OUTPUT_API} -e staging
run.api.dev:
	reflex -c ./reflex.conf
clean:
	rm -rf ${OUTPUT_DIR}
debug: clean build.api.dev dlv.api
dlv.api:
	dlv exec ${OUTPUT_API} --listen=${DEBUG_PORT_API} --headless=true --api-version=2 --accept-multiclient
