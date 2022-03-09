#!/bin/bash
version=`grep version main.go  | awk '{print $3}'`
docker build -t registry.cn-hangzhou.aliyuncs.com/my_docker_images/httpserver:$version -f DockerfileMultiSegmentBuild  .
docker push registry.cn-hangzhou.aliyuncs.com/my_docker_images/httpserver:$version