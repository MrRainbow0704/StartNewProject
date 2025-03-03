#!/usr/bin/env bash

echo "Killing gunicorn..."
cat /var/run/gunicorn/dev.pid | xargs kill -9
sleep 5
tail /var/log/gunicorn/dev.log -n 50
echo "Gunicorn killed!"