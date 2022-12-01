#!/bin/sh
if [ -z "${AWS_LAMBDA_RUNTIME_API}" ]; then
#   exec /usr/bin/aws-lambda-rie "$@"
  exec /usr/local/bin/aws-lambda-rie "$@"
else
  exec "$@"
fi  