#!/bin/bash

echo -n $(head -c 16 /dev/urandom | shasum | cut -d " " -f 1)
