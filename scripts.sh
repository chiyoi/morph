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
    echo "$0 log"
    echo "    Track container log."
    echo "$0 up"
    echo "    Run in docker."
    echo "$0 stop"
    echo "    Stop and clear running container."
    echo "$0 update"
    echo "    Stop running container, Pull, Build and Up new version."
}

ARTIFACT=morph

dev_env() {
    export VERSION="dev"
    export ADDR=":12380"
    export ENV="dev"
    export ENDPOINT_AZURE_COSMOS="https://neko03cosmos.documents.azure.com:443/"
    export ENDPOINT_AZURE_BLOB="https://neko03storage.blob.core.windows.net/"
    export BLOB_CONTAINER_CERT_CACHE="neko0001"
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
    sudo docker build -t chiyoi/$ARTIFACT .
}

up() {
    sudo docker run -d --network=host --restart=on-failure:5 --name=$ARTIFACT chiyoi/$ARTIFACT
}

log() {
    sudo docker logs -f $ARTIFACT
}

stop() {
    sudo docker stop $ARTIFACT && sudo docker rm $ARTIFACT
}

update() {
    stop 2>/dev/null
    pull && build && up
}

if test -z "$1" -o -n "$(echo "$1" | grep -Ex '\-{0,2}h(elp)?')"; then
usage
exit
fi

case "$1" in
tidy|run|pull|build|up|log|stop|update);;
*)
usage
exit 1
;;
esac

$@

