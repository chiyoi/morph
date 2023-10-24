#!/bin/sh
cd $(dirname $(realpath $0)) || return
usage() {
    pwd
    echo "Scripts:"
    echo "$0 tidy"
    echo "    Tidy go module."
    echo "$0 run"
    echo "    Run the main package."
    echo "$0 pull"
    echo "    Pull from git origin."
    echo "$0 build"
    echo "    Build docker image."
    echo "$0 up"
    echo "    Run in docker."
    echo "$0 update"
    echo "    Pull, Build and Up."
}

dev_env() {
    export ADDR=":12380"
    export ENV="dev"
    export ENDPOINT_AZURE_COSMOS="https://neko03cosmos.documents.azure.com:443/"
    export DATABASE="neko0001"
}

tidy() {
    go mod tidy
}

run() {
    dev_env
    go run .
}

pull() {
    git pull
}

build() {
    docker build -t chiyoi/morph .
}

up() {
    docker run -d --network=host --restart=on-failure:5 --name=morph chiyoi/morph
}

update() {
    pull && build && up
}

if test -z "$1" -o -n "$(echo "$1" | grep -Ex '\-{0,2}h(elp)?')"; then
usage
exit
fi

case "$1" in
tidy|run|pull|build|up|update);;
*)
usage
exit 1
;;
esac

$@

