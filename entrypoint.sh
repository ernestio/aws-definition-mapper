#!/usr/bin/env sh

echo "Waiting for NATS"
while ! echo exit | nc nats 4222; do sleep 1; done

echo "Waiting for Postgres"
while ! echo exit | nc postgres 5432; do sleep 1; done

echo "Starting aws-definition-mapper"
/go/bin/aws-definition-mapper
