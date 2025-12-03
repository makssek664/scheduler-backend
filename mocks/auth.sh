#!/bin/sh 
curl -X POST -d '{"name": "newuser"}' localhost:8080/auth
echo "\n"
