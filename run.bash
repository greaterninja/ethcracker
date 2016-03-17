#!/bin/sh

go run src/ethcracker.go -pk ~/test/pk.txt -t ~/test/templates.txt -threads 4  -min_len 1
