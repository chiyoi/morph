#!/bin/sh
cd $(dirname $(realpath $0)) || return
usage() {
    pwd
    echo "Scripts:"
    echo "$0 tidy"
    echo "    Go mod tidy."
    echo "$0 run"
    echo "    Run the main package."
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

if test -z "$1" -o -n "$(echo "$1" | grep -Ex '\-{0,2}h(elp)?')"; then
usage
exit
fi

case "$1" in
tidy|run);;
*)
usage
exit 1
;;
esac

$@

