#!/usr/bin/env bash

i=0

while [ $i -le 30 ]; do
	printf '%s line\n' "$i"
	i=$((i + 1))
	sleep .1
done
