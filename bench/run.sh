#!/usr/bin/env bash

wrk -t10 -c300 -d10s --script $1 http://localhost:8080/token
