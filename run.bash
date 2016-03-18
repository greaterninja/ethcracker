#!/bin/sh

go run src/ethcracker.go -pk ~/test/s.txt -t ~/test/ethcracker-pwd2.txt -threads 4  -min_len 1 -max_len 40 -v 1 -start_from 0 -re 1
#go run src/ethcracker.go -pk ~/test/ethwallet-q.json -t ~/test/pattern.txt -presale -threads 4  -min_len 1 -v 1 -start_from 0
