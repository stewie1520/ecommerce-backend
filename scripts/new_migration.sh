#! /bin/sh
# rebuild prog if necessary
make clean
make build.api
# run prog with some arguments
./bin/api migrate create "$@"
