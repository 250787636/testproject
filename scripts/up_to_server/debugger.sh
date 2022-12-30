HOST="tcp://172.16.102.58:2375"
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
debug backend-ssp-user  #user端 docker debug端口8001
debug backend-ssp-admin #admin端 docker debug端口8010