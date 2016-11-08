#!/bin/bash

test -d logs || mkdir logs

./fluentd_conf & 

