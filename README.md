# ethcracker
Ethereum password cracker


# Usage 

    ethcracker -pk ~/test/pk.txt -t ~/test/templates.txt

    -pk path to the private key file
    -t  path to the template file
    -presale  for cracking prelase JSSON file

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

# Donation

If this program helps you to restore the password, please donate some ETH to the address:

 0x281694Fabfdd9735e01bB59942B18c469b6e3df6
 
 Thank you