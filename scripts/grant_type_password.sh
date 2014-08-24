#!/usr/bin/env bash

curl -v --user client_id:client_secret --data "grant_type=password&username=email@example.com&password=password" http://localhost:8080/token
