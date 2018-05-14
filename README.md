# Encryption Playground

* See `private-public-key-and-diffie-hellmann.ods` for illustration.  

* Open the file with Open Office Calc or with MS Excel.

* Primitive setup to find a divider and two corresponding numbers as public and private key.

* Uses golang math/big to process huge numbers.

## Implementation in Go

* `inp.json` contains two helper primes and some plain text to encrypt/decrypt.

* main() computes a product of two primes serving as divider 

* It chooses a random number as public key (does it have to be prime?)

* It then searched for a corresponding private key, so that 

    (pub*priv) % modified_divider == 1

## Encrypt - decrypt

* The plain text will be encrypted using the public key

* The cipher text will be decrypted using the private key

* Results are saved to out.json

## Beware of amateur details

* The search for a related private key is deterministic  
 and does not scale.