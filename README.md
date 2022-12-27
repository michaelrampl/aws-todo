# Todo API

This repository contains 2 different binaries.

## Standalone

This binary implements the fiber framework and can be used as a standalone version for servers.

## Lambda

This binary is intended to be used as AWS Lambda function behind an API Gateway.

## Structure

Both binaries access the same handler but inject a different database (all based on a common interface). At the time being, either mongodb or dynamodb are possible backends.