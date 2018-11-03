#!/usr/bin/env bash

sshpass -p "123123" ssh -t -p 9844 reza@188.40.164.251 "docker logs -t -f calendar"
