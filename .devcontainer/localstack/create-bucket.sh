#!/bin/sh
# Runs when LocalStack is ready. With persisted data the bucket survives restarts.
# awslocal defaults to port 4566; LocalStack listens on GATEWAY_LISTEN (see docker-compose.yml).
export AWS_ENDPOINT_URL=http://localhost:3171
awslocal s3api head-bucket --bucket deliver 2>/dev/null || awslocal s3 mb s3://deliver
