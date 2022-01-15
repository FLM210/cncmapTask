#!/bin/bash


docker run -itd  --name httpserver hahtangtang/httpserver:v1 

nsenter -t `docker inspect -f {{.State.Pid}} httpserver` -n