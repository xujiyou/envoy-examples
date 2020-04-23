#!/bin/sh
./code/main_v2 &
envoy -c /etc/envoy-config-v2.yaml --service-cluster service-v2
