#!/bin/zsh
scripts=$0
cd $(dirname $(realpath $scripts)) || return
usage () {
    pwd
    echo "Scripts:"
    echo "$scripts help"
    echo "    Show this help message."
    echo "$scripts tidy"
    echo "    Tidy go module."
    echo "$scripts run"
    echo "    Run the main package."
    echo "$scripts build"
    echo "    Build docker image."
    echo "$scripts logs"
    echo "    Track container log."
    echo "$scripts up"
    echo "    Run in docker."
    echo "$scripts stop"
    echo "    Stop and clear running container."
    echo "$scripts update"
    echo "    Stop running container, Pull, Build and Up new version."
}

export ENV="dev"
ARTIFACT=morph

help () {
    usage
}

tidy () {
    go mod tidy
}

run () {
    go run .
}

build () {
    sudo docker build -t chiyoi/$ARTIFACT .
}

up () {
    sudo docker run -d --network=host --restart=on-failure:5 --name=$ARTIFACT chiyoi/$ARTIFACT
}

logs () {
    sudo docker logs -f $ARTIFACT
}

stop () {
    sudo docker stop $ARTIFACT && sudo docker rm $ARTIFACT
}

update () {
    git pull && build || return
    stop 2>/dev/null
    up
}

case "$1" in
""|-h|-help|--help)
usage
exit
;;
help|tidy|run|build|up|logs|stop|update)
$@
;;
*)
usage
exit 1
;;
esac
