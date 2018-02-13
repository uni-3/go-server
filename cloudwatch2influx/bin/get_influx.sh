#!/bin/bash


curl -G -XGET 'http://localhost:8086/query?pretty=true' --data-urlencode "db=ggg" --data-urlencode 'q=SELECT "idle","system","user" FROM cpu_usage'
