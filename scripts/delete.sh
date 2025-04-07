#!/usr/bin/env bash
echo '{"name":"sackbuoy-mc", "namespace":"games"}' | \
  http DELETE 'http://localhost:8080/servers/delete'

