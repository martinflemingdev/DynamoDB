#!/usr/bin/bash
aws dynamodb describe-table --table-name Music | grep TableStatus