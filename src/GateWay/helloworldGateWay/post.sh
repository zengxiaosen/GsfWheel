#!/bin/bash
#curl -X POST localhost:8083/hello -d '{"name": "rephus"}'
curl -X POST -k http://localhost:8083/api/hello -d '{"name": " world"}'
