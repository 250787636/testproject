#!/usr/bin/env zsh

export GOOS=linux
export GOARCH=amd64

# 服务器地址
HOST="tcp://172.16.102.58:2375"

#开启远程调试
enableDebugger=true

function debug() {
#    echo "kill ./program"
#    docker -H ${HOST} exec "$1" sh -c "ps -ef | grep program | grep -v grep | awk '{print \$2}'| xargs kill -9"
    retStr=$(docker -H ${HOST} exec "$1" sh -c "ls /data/bin/program|grep dlv")
    if [ "$retStr" != 'dlv' ];then
        echo "copy debugger tool from ./dlv to /data/bin/program"
        docker -H ${HOST} cp ./dlv "$1":/data/bin/program/dlv
    fi
    pid=$(docker -H ${HOST} exec "$1" sh -c "ps -ef | grep program | grep -v grep | awk '{print \$2}'")
    echo "$pid"
    docker -H ${HOST} exec "$1" sh -c "chmod 777 -R  /data/bin/program/dlv; nohup /data/bin/program/dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient attach ${pid} &"
    return
}

echo "build user"
go build -o program_user ../../main.go
echo "build admin"
go build -o program_admin ../../adminbackend/main.go
echo "build finish"
echo "docker cp"
docker -H ${HOST} cp ./program_user backend-ssp-user:/data/bin/program/program
docker -H ${HOST} cp ./program_admin backend-ssp-admin:/data/bin/program/program
docker -H ${HOST} cp ../../sconf backend-ssp-user:/data/bin/program
docker -H ${HOST} cp ../../adminbackend/sconf backend-ssp-admin:/data/bin/program
echo "cp finish"

docker -H ${HOST} restart backend-ssp-admin
docker -H ${HOST} restart backend-ssp-user

if [ $enableDebugger = true ];then
    echo "start debug backend-ssp-user" #user端 docker debug端口8001
    debug backend-ssp-user
    echo "start debug backend-ssp-admin" #admin端 docker debug端口8010
    debug backend-ssp-admin
fi
rm  ./program_user
rm ./program_admin
