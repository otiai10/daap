#!/bin/sh

MACHINE=daaptest

printf "Checking docker-machine CLI: "
if [ -z "`which docker-machine`" ]; then
  echo "NG! 'docker-machine' command not found"
  exit 1
fi
printf "ok\n"

printf "Preparing test machine: "
if [ -z "`docker-machine ls -q | grep ${MACHINE}`" ]; then
  printf "\n\tCreating test machine '${MACHINE}': "
  docker-machine create ${MACHINE}
  printf "ok\n"
else
  printf "\n\tTest machine '${MACHINE}' already exits: ok\n"
fi

printf "Preparing test data: "
docker-machine scp -r ./tests/testdata ${MACHINE}:/home/docker/data > /dev/null
printf "ok\n"

printf "Installing env variables: "
eval $(docker-machine env ${MACHINE})
printf "ok\n"

printf "Installing golang dependencies: "
dep ensure
printf "ok\n"

printf "Running golang test:\n"
go test -v .
