#!/usr/bin/env bash
echo '{"name":"sackbuoy-mc", "namespace":"games", "gameType":"minecraft"}' | \
  http POST 'http://localhost:8080/servers/create'

