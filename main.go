package main

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/zew/util"
)

type Set struct {
	Pri1 *big.Int `json:"prime_1,omitempty"` // A helper prime
	Pri2 *big.Int `json:"prime_2,omitempty"` // A helper prime
	P    *big.Int `json:"p,omitempty"`       // Our divider, product of helper primes 1 and 2

	// A helper modification of P, to compute matching private and public exponent.
	// Computed as prime_1 -1 * prime_2 - 1
	PDecr *big.Int `json:"p_decremented,omitempty"`

	ExpPub *big.Int `json:"exp_private,omitempty"`
	ExpPrv *big.Int `json:"exp_public,omitempty"`
	// (ExpPrv*ExpPub) %  PDecr == 1  !

	Plain         string     `json:"plain_text,omitempty"`
	PlainInts     []*big.Int `json:"plain_ints,omitempty"`
	Cypher        []*big.Int `json:"cypher,omitempty"`
	DecryptedInts []*big.Int `json:"decr_ints,omitempty"`
	Decrpyted     string     `json:"decr_text,omitempty"`
}

func (s *Set) Init() bool {
	s.P = big.NewInt(0)
	s.P.Mul(s.Pri1, s.Pri2)

	s.PDecr = big.NewInt(0)
	h1, h2 := big.NewInt(0), big.NewInt(0)
	h1.Sub(s.Pri1, big.NewInt(1))
	h2.Sub(s.Pri2, big.NewInt(1))
	s.PDecr.Mul(h1, h2)
	log.Printf("Prime1     = %3v   -  Prime2     = %3v  | %4v", s.Pri1, s.Pri2, s.P)
	log.Printf("Prime1 - 1 = %3v   -  Prime2 - 1 = %3v  | %4v", h1, h2, s.PDecr)

	log.Printf("----------")
	//
	// Random seed for public exponent
	P_decr := s.PDecr.Int64()
	expPub := rand.Int63n(P_decr-2) + int64(2)
	s.ExpPub = big.NewInt(expPub)
	log.Printf("Exp Public is %3v", s.ExpPub)

	//
	// Find an appropriate private key with
	// REST( ExpPrv * ExpPub ; PDecr ) == 1  !
	maxTries := s.P.Int64()
	fmtStr := "%3v * %3v MOD %3v = %3v"

	for i := int64(2); i < maxTries; i++ {

		// Do not consider non-primes
		iBig := big.NewInt(i)
		isPrime := iBig.ProbablyPrime(1)
		if !isPrime {
			// log.Printf("skipping for non prime: %v", iBig)
			continue
		}

		// private key should be unequal public key
		if iBig.Cmp(s.ExpPub) == 0 {
			log.Printf("skipping %v", iBig)
			continue
		}

		msg := ""
		re := (expPub * i) % P_decr
		if re == 1 {
			s.ExpPrv = iBig
			msg = fmt.Sprintf("| found a second number %3v, so that remainder is 1", i)
			log.Printf(fmtStr+"  %v", expPub, i, P_decr, re, msg)
			log.Printf("(%v*%v)%%%v = %v", expPub, i, s.P.Int64(), (expPub*i)%s.P.Int64())
			return true
		}
		log.Printf(fmtStr, expPub, i, P_decr, re)
	}
	return false
}

func (s *Set) Encrypt(plain []*big.Int) (cypher []*big.Int) {
	// =REST(plain^ExpPub;P)
	for _, v := range plain {
		tmp := big.NewInt(0)
		tmp.Exp(v, s.ExpPub, s.P)
		cypher = append(cypher, tmp)
	}
	return
}
func (s *Set) Decrypt(cypher []*big.Int) (plain []*big.Int) {
	// =REST(plain^ExpPub;P)
	for _, v := range cypher {
		tmp := big.NewInt(0)
		tmp.Exp(v, s.ExpPrv, s.P)
		plain = append(plain, tmp)
	}
	return
}

func main() {

	log.SetFlags(log.Lshortfile)
	rand.Seed(int64(time.Now().Nanosecond()))

	s1 := Load()
	log.Printf("%v", util.IndentedDump(s1))
	// s1.Pri2 = big.NewInt(11)

	for {
		found := s1.Init()
		if found {
			break
		} else {
			log.Printf("No priv key found")
		}
	}
	log.Printf("%v", util.IndentedDump(s1))

	for _, char := range s1.Plain {
		s1.PlainInts = append(s1.PlainInts, big.NewInt(int64(char)))
	}
	s1.Cypher = s1.Encrypt(s1.PlainInts)
	// log.Printf("%v", util.IndentedDump(s1))
	// log.Printf("%v", util.IndentedDump(s1.Cypher))

	s1.DecryptedInts = s1.Decrypt(s1.Cypher)
	x := []byte{}
	for _, anInt := range s1.DecryptedInts {
		x = append(x, byte(anInt.Int64()))
	}
	s1.Decrpyted = string(x)

	s1.Save()
	log.Printf("%v", util.IndentedDump(s1))

}
