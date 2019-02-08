#!/bin/bash

set -eu -o pipefail


for f in examples/*; do
	if [[ -d $f ]]; then
		go build ./$f;
	fi;
done
