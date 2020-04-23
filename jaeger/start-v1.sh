#!/bin/sh
./code/main_v1 &
envoy -c /etc/envoy-config-v1.yaml --service-cluster service-v1
