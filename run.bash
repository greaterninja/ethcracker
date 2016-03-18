#!/bin/sh

go run src/ethcracker.go -pk ~/test/pk.txt -t ~/test/templates.txt -threads 4  -min_len 1 -v 1 -start_from 0 -re 7
#go run src/ethcracker.go -pk ~/test/ethwallet-q.json -t ~/test/pattern.txt -presale -threads 4  -min_len 1 -v 1 -start_from 0
