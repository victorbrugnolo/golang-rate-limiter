#/bin/bash

for i in {1..100} ; do
  echo ' http://localhost:8080/ping'
done | xargs curl -s