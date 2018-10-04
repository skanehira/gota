#!/bin/bash

# build image
docker build -t gota .

# remove build image
docker rmi $(docker images --filter "dangling=true" -aq)

# push image to dockerr hub
#docker push skanehira/gota

