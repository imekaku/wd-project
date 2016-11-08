#!/bin/bash

if ps -ef | grep fluentd_conf | grep -v grep
then
    killall fluentd_conf
fi
