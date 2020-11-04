#!/bin/sh

for f in bin/ofc*; do shasum -a 256 $f > $f.sha256; done
