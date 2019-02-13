#!/bin/bash

# docker_compose.yml from openfaas/faas tag 0.11.1
cd ./tmp/faas && docker stack deploy func -c ./docker-compose.yml