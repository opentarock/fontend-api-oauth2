#!/usr/bin/env bash

curl --user client_id:client_secret --data "grant_type=refresh_token&refresh_token=$1" http://localhost:8080/token
