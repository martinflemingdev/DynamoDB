#!/usr/bin/bash

# The default behavior for DynamoDB is eventually consistent reads. The consistent-read parameter is used below to demonstrate strongly consistent reads.

aws dynamodb get-item --consistent-read \
    --table-name Music \
    --key '{ "Artist": {"S": "Acme Band"}, "SongTitle": {"S": "Happy Day"}}'