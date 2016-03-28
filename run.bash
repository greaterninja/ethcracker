#!/bin/sh

go run src/ethcracker.go -pk ~/test/pk.txt -t ~/test/ethcracker-pwd3.txt -keep_order -threads 4  -min_len 1 -max_len 40 -v 1 -start_from 0 -dump ~/test/v.txt -re 1  
#go run src/ethcracker.go -pk ~/test/ethwallet-q.json -t ~/test/pattern.txt -presale -threads 4  -min_len 1 -v 1 -start_from 0
