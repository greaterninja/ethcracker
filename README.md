# Ethereum password cracker

If you forgot your password for Ethereum private key or presale file, but think you still remember 
the list of possible substrings from wich you constracted the passwerd, then try this progam to 
routinly go through all the possible combinations and find the working password.

PLEASE DO NOT TRUST ANYONE TO COMPILE THE PROGRAM FOR YOU. ALWAYS USE THE SOURCE CODE DOWNLOADED FROM 
THE GITHUB. THIS WAY YOU CAN BE SURE, THE PROGRAM DOES NOT HAVE ANY MALICIOUS CODE !!!

# Usage 

    ethcracker -pk ~/test/pk.txt -t ~/test/templates.txt

    -pk path to the private key file
    -t  path to the template file
    -presale  for cracking prelase JSSON file
    -threads Number of threads
    -v Verbosity ( 0, 1, 2 )
    -start_from Skip first N combinations

# Template file format

Every line contains the possible variants of the substring. For example file:

    a1 a2 a3
    b

will generate all those combinations

    a1
    a2
    a3
    b
    a1b
    ba1
    a2b
    ba2
    a3b
    ba3

# Installing

Install Go Language


    git clone https://github.com/lexansoft/ethcracker
    cd ethcracker
    
    go get github.com/ethereum/go-ethereum/common
    go get github.com/ethereum/go-ethereum/crypto/ecies
    go get github.com/ethereum/go-ethereum/crypto/randentropy
    go get github.com/ethereum/go-ethereum/crypto/secp256k1
    go get github.com/ethereum/go-ethereum/crypto/sha3
    go get github.com/ethereum/go-ethereum/rlp
    go get github.com/pborman/uuid
    go get golang.org/x/crypto/pbkdf2
    go get golang.org/x/crypto/ripemd160
    go get golang.org/x/crypto/scrypt
    
    go run src/ethcracker.go -pk PATH_TO_FILE -t PATH_TO_TAMPLATE_FILE -threads 4 
    


# Donation

If this program helps you to restore the password, please donate some ETH to the address:

 0x281694Fabfdd9735e01bB59942B18c469b6e3df6
 
 Thank you